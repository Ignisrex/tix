package booking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/ignisrex/tix/booking/internal/payment"
	"github.com/ignisrex/tix/booking/internal/redis"
	"github.com/ignisrex/tix/booking/types"
)

type Service struct {
	repo        *Repo
	redisClient *redis.Client
}

func NewService(repo *Repo, redisClient *redis.Client) *Service {
	return &Service{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *Service) Reserve(ctx context.Context, id uuid.UUID) (*types.ReserveResponse, int, error) {
	
	status, err := s.repo.GetTicketStatus(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &types.ReserveResponse{
				Success:  false,
				Message:  "ticket not found",
				TicketID: id,
			}, http.StatusNotFound, fmt.Errorf("ticket not found: %w", err)
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to get ticket status: %w", err)
	}
	if status == "sold" {
		return &types.ReserveResponse{
			Success:  false,
			Message:  "ticket is already sold",
			TicketID: id,
		}, http.StatusGone, errors.New("ticket already sold")
	}

	// Check if ticket is already reserved in Redis
	reserved, err := s.redisClient.IsReserved(ctx, id)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to check reservation status: %w", err)
	}
	if reserved {
		return &types.ReserveResponse{
			Success:  false,
			Message:  "ticket is already reserved",
			TicketID: id,
		}, http.StatusConflict, errors.New("ticket already reserved")
	}

	// Try to reserve the ticket (atomic operation with SETNX)
	set, err := s.redisClient.ReserveTicket(ctx, id)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to reserve ticket: %w", err)
	}
	if !set {
		// Another request reserved it between check and set
		return &types.ReserveResponse{
			Success:  false,
			Message:  "ticket is already reserved",
			TicketID: id,
		}, http.StatusConflict, errors.New("ticket already reserved")
	}

	return &types.ReserveResponse{
		Success:  true,
		Message:  "ticket reserved successfully",
		TicketID: id,
	}, http.StatusOK, nil
}

func (s *Service) Purchase(ctx context.Context, id uuid.UUID) (*types.PurchaseResponse, int, error) {
	/*
	    If user was integrated into the system, 
		we would check to make sure the user who is purchasing the ticket 
		is the same user who reserved the ticket.
	*/
	
	ticket, err := s.repo.GetTicketWithPrice(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &types.PurchaseResponse{
				Success:  false,
				Message:  "ticket not found",
				TicketID: id,
				Total:    0,
			}, http.StatusNotFound, fmt.Errorf("ticket not found: %w", err)
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to get ticket details: %w", err)
	}
	
	if ticket.Status == "sold" {
		return &types.PurchaseResponse{
			Success:  false,
			Message:  "ticket is already sold",
			TicketID: id,
			Total:    0,
		}, http.StatusGone, errors.New("ticket already sold")
	}


	ticketInfo := payment.TicketInfo{
		ID:         ticket.ID,
		PriceCents: ticket.PriceCents,
	}
	
	if err := payment.ProcessPayment(ctx, ticketInfo); err != nil {
		return &types.PurchaseResponse{
			Success:  false,
			Message:  fmt.Sprintf("payment failed: %v", err),
			TicketID:  id,
			Total:    ticket.PriceCents,
		}, http.StatusPaymentRequired, err
	}

	if err := s.repo.PurchaseTicket(ctx, id); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to update ticket status: %w", err)
	}

	if err := s.redisClient.ReleaseTicket(ctx, id); err != nil {
		log.Printf("failed to release ticket: %v", err)
	}

	return &types.PurchaseResponse{
		Success:  true,
		Message:  "purchase completed successfully",
		TicketID: id,
		Total:    ticket.PriceCents,
	}, http.StatusOK, nil
}


