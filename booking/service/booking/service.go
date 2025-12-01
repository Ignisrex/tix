package booking

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

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

// ReserveTickets attempts to reserve multiple tickets atomically.
// If any ticket fails to reserve, all reservations are rolled back.
func (s *Service) ReserveTickets(ctx context.Context, ticketIDs []uuid.UUID) (*types.ReserveResponse, int, error) {
	if len(ticketIDs) == 0 {
		return &types.ReserveResponse{
			Success:   false,
			Message:   "no tickets provided",
			TicketIDs: []uuid.UUID{},
		}, http.StatusBadRequest, errors.New("no tickets provided")
	}

	// Validate all tickets exist and are available
	for _, id := range ticketIDs {
		status, err := s.repo.GetTicketStatus(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return &types.ReserveResponse{
					Success:   false,
					Message:   fmt.Sprintf("ticket %s not found", id),
					TicketIDs: []uuid.UUID{},
				}, http.StatusNotFound, fmt.Errorf("ticket %s not found: %w", id, err)
			}
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to get ticket status: %w", err)
		}
		if status == "sold" {
			return &types.ReserveResponse{
				Success:   false,
				Message:   fmt.Sprintf("ticket %s is already sold", id),
				TicketIDs: []uuid.UUID{},
			}, http.StatusGone, fmt.Errorf("ticket %s already sold", id)
		}

		// Check if ticket is already reserved
		reserved, err := s.redisClient.IsReserved(ctx, id)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to check reservation status: %w", err)
		}
		if reserved {
			return &types.ReserveResponse{
				Success:   false,
				Message:   fmt.Sprintf("ticket %s is already reserved", id),
				TicketIDs: []uuid.UUID{},
			}, http.StatusConflict, fmt.Errorf("ticket %s already reserved", id)
		}
	}

	// Attempt to reserve all tickets atomically
	results, err := s.redisClient.ReserveTickets(ctx, ticketIDs)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to reserve tickets: %w", err)
	}

	// Check if all reservations succeeded
	var reservedIDs []uuid.UUID
	var failedIDs []uuid.UUID
	for ticketID, success := range results {
		if success {
			reservedIDs = append(reservedIDs, ticketID)
		} else {
			failedIDs = append(failedIDs, ticketID)
		}
	}

	// If any failed, rollback all successful reservations
	if len(failedIDs) > 0 {
		if len(reservedIDs) > 0 {
			_ = s.redisClient.ReleaseTickets(ctx, reservedIDs)
		}
		return &types.ReserveResponse{
			Success:   false,
			Message:   "one or more tickets could not be reserved",
			TicketIDs: []uuid.UUID{},
		}, http.StatusConflict, fmt.Errorf("failed to reserve tickets: %v", failedIDs)
	}

	return &types.ReserveResponse{
		Success:   true,
		Message:   "tickets reserved successfully",
		TicketIDs: reservedIDs,
	}, http.StatusOK, nil
}

// PurchaseTickets attempts to purchase multiple tickets atomically.
// If any ticket fails, all operations are rolled back.
func (s *Service) PurchaseTickets(ctx context.Context, ticketIDs []uuid.UUID) (*types.PurchaseResponse, int, error) {
	if len(ticketIDs) == 0 {
		return &types.PurchaseResponse{
			Success:    false,
			Message:    "no tickets provided",
			TicketIDs:  []uuid.UUID{},
			Total:      0,
			PurchaseID: uuid.Nil,
		}, http.StatusBadRequest, errors.New("no tickets provided")
	}

	// Get all ticket details
	tickets := make([]struct {
		ID         uuid.UUID
		PriceCents int32
		Status     string
	}, 0, len(ticketIDs))
	
	totalCents := int32(0)
	for _, id := range ticketIDs {
		ticket, err := s.repo.GetTicketWithPrice(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return &types.PurchaseResponse{
					Success:    false,
					Message:    fmt.Sprintf("ticket %s not found", id),
					TicketIDs:  []uuid.UUID{},
					Total:      0,
					PurchaseID: uuid.Nil,
				}, http.StatusNotFound, fmt.Errorf("ticket %s not found: %w", id, err)
			}
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to get ticket details: %w", err)
		}

		if ticket.Status == "sold" {
			return &types.PurchaseResponse{
				Success:    false,
				Message:    fmt.Sprintf("ticket %s is already sold", id),
				TicketIDs:  []uuid.UUID{},
				Total:      0,
				PurchaseID: uuid.Nil,
			}, http.StatusGone, fmt.Errorf("ticket %s already sold", id)
		}

		tickets = append(tickets, struct {
			ID         uuid.UUID
			PriceCents int32
			Status     string
		}{
			ID:         ticket.ID,
			PriceCents: ticket.PriceCents,
			Status:     string(ticket.Status),
		})
		totalCents += ticket.PriceCents
	}

	// Process payment for all tickets
	ticketInfos := make([]payment.TicketInfo, 0, len(tickets))
	for _, ticket := range tickets {
		ticketInfos = append(ticketInfos, payment.TicketInfo{
			ID:         ticket.ID,
			PriceCents: ticket.PriceCents,
		})
	}

	if err := payment.ProcessPayment(ctx, ticketInfos...); err != nil {
		return &types.PurchaseResponse{
			Success:    false,
			Message:    fmt.Sprintf("payment failed: %v", err),
			TicketIDs:  []uuid.UUID{},
			Total:      totalCents,
			PurchaseID: uuid.Nil,
		}, http.StatusPaymentRequired, err
	}

	// Purchase all tickets in a transaction
	purchaseID, err := s.repo.PurchaseTickets(ctx, ticketIDs, totalCents)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to purchase tickets: %w", err)
	}

	// Release all reservations
	if err := s.redisClient.ReleaseTickets(ctx, ticketIDs); err != nil {
		log.Printf("failed to release tickets: %v", err)
	}

	return &types.PurchaseResponse{
		Success:   true,
		Message:   "purchase completed successfully",
		TicketIDs: ticketIDs,
		Total:     totalCents,
		PurchaseID: purchaseID,
	}, http.StatusOK, nil
}

// GetPurchaseDetails retrieves purchase details including all tickets
func (s *Service) GetPurchaseDetails(ctx context.Context, purchaseID uuid.UUID) (*types.PurchaseDetailsResponse, int, error) {
	details, err := s.repo.GetPurchaseDetails(ctx, purchaseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, http.StatusNotFound, fmt.Errorf("purchase not found: %w", err)
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to get purchase details: %w", err)
	}

	// Parse tickets JSON
	var ticketDetails []types.PurchaseTicketDetail
	if err := json.Unmarshal(details.Tickets, &ticketDetails); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to parse ticket details: %w", err)
	}

	return &types.PurchaseDetailsResponse{
		PurchaseID:        details.PurchaseID,
		TotalCents:        details.TotalCents,
		PurchaseCreatedAt: details.PurchaseCreatedAt.Format(time.RFC3339),
		Tickets:           ticketDetails,
	}, http.StatusOK, nil
}


