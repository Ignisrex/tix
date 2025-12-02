package booking

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

// Domain-level error markers used by handlers to map to HTTP responses.
var (
	ErrTicketNotFound   = errors.New("ticket not found")
	ErrTicketSold       = errors.New("ticket sold")
	ErrTicketReserved   = errors.New("ticket reserved")
	ErrPaymentFailed    = errors.New("payment failed")
	ErrPurchaseNotFound = errors.New("purchase not found")
)

func NewService(repo *Repo, redisClient *redis.Client) *Service {
	return &Service{
		repo:        repo,
		redisClient: redisClient,
	}
}

// ReserveTickets validates that all tickets exist and are available, and then
// attempts to reserve them atomically in Redis.
// On success it returns the reserved ticket IDs; on failure it returns a
// domain error (e.g. ErrTicketNotFound, ErrTicketSold).
func (s *Service) ReserveTickets(ctx context.Context, ticketIDs []uuid.UUID) ([]uuid.UUID, error) {
	// Validate all tickets exist and are available
	tickets, err := s.repo.GetTicketsWithPrice(ctx, ticketIDs)
	if err != nil {
		log.Printf("ReserveTickets: failed to get tickets with price: %v", err)
		return nil, fmt.Errorf("failed to get tickets with price: %w", err)
	}

	if len(tickets) != len(ticketIDs) {
		log.Printf("ReserveTickets: some tickets not found (requested=%d, found=%d)", len(ticketIDs), len(tickets))
		return nil, fmt.Errorf("%w: some tickets not found", ErrTicketNotFound)
	}

	for _, ticket := range tickets {
		//check if ticket is sold
		if ticket.Status == "sold" {
			log.Printf("ReserveTickets: ticket %s already sold", ticket.ID)
			return nil, fmt.Errorf("%w: ticket %s already sold", ErrTicketSold, ticket.ID)
		}
	}

	// Attempt to reserve all tickets atomically
	if err := s.redisClient.ReserveTickets(ctx, ticketIDs); err != nil {
		log.Printf("ReserveTickets: failed to reserve tickets in redis: %v", err)
		return nil, fmt.Errorf("failed to reserve tickets: %w", err)
	}

	return ticketIDs, nil
}

// PurchaseTickets attempts to purchase multiple tickets atomically.
// If any ticket fails, all operations are rolled back and tickets are released
// It returns the purchase ID and total cents on success.
// On failure it returns a domain error (e.g. ErrTicketNotFound, ErrPaymentFailed).
func (s *Service) PurchaseTickets(ctx context.Context, ticketIDs []uuid.UUID) (uuid.UUID, int32, error) {
	// Refresh the lock TTL for each ticket to 10 minutes while processing payment
	ok, err := s.redisClient.RefreshTickets(ctx, ticketIDs, 10*time.Minute)
	if err != nil {
		log.Printf("failed to refresh ticket locks before purchase: %v", err)
	}

	if !ok {
		log.Printf("PurchaseTickets: one or more tickets are not reserved at purchase time")
		return uuid.Nil, 0, fmt.Errorf("%w: one or more tickets are not reserved", ErrTicketReserved)
	}

	tickets, err := s.repo.GetTicketsWithPrice(ctx, ticketIDs)
	if err != nil {
		log.Printf("PurchaseTickets: failed to get ticket details: %v", err)
		return uuid.Nil, 0, fmt.Errorf("failed to get ticket details: %w", err)
	}

	if len(tickets) != len(ticketIDs) {
		log.Printf("PurchaseTickets: some tickets not found (requested=%d, found=%d)", len(ticketIDs), len(tickets))
		return uuid.Nil, 0, fmt.Errorf("%w: some tickets not found", ErrTicketNotFound)
	}

	totalCents, err := payment.ProcessPayment(ctx, tickets...)
	if err != nil {
		log.Printf("PurchaseTickets: payment failed: %v", err)
		return uuid.Nil, 0, fmt.Errorf("%w: %v", ErrPaymentFailed, err)
	}

	// Purchase all tickets in a transaction
	purchaseID, err := s.repo.PurchaseTickets(ctx, ticketIDs, int32(totalCents))
	if err != nil {
		log.Printf("PurchaseTickets: failed to purchase tickets in db: %v", err)
		return uuid.Nil, 0, fmt.Errorf("failed to purchase tickets: %w", err)
	}

	// Release all reservations
	if err := s.redisClient.ReleaseTickets(ctx, ticketIDs); err != nil {
		log.Printf("failed to release tickets: %v", err)
	}

	return purchaseID, int32(totalCents), nil
}

// GetPurchaseDetails retrieves purchase details including all tickets
func (s *Service) GetPurchaseDetails(ctx context.Context, purchaseID uuid.UUID) (*types.PurchaseDetailsResponse, error) {
	details, err := s.repo.GetPurchaseDetails(ctx, purchaseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("GetPurchaseDetails: purchase %s not found", purchaseID)
			return nil, fmt.Errorf("%w: %v", ErrPurchaseNotFound, err)
		}
		log.Printf("GetPurchaseDetails: failed to get purchase details from db: %v", err)
		return nil, fmt.Errorf("failed to get purchase details: %w", err)
	}

	var ticketDetails []types.PurchaseTicketDetail
	if err := json.Unmarshal(details.Tickets, &ticketDetails); err != nil {
		log.Printf("GetPurchaseDetails: failed to parse ticket details json: %v", err)
		return nil, fmt.Errorf("failed to parse ticket details: %w", err)
	}

	resp := &types.PurchaseDetailsResponse{
		PurchaseID:        details.PurchaseID,
		TotalCents:        details.TotalCents,
		PurchaseCreatedAt: details.PurchaseCreatedAt.Format(time.RFC3339),
		Tickets:           ticketDetails,
	}

	return resp, nil
}

// CheckTicketLocks checks the reservation status for multiple tickets.
// Returns a map of ticketID -> is_reserved (true if reserved, false if available).
func (s *Service) CheckTicketLocks(ctx context.Context, ticketIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	return s.redisClient.AreReserved(ctx, ticketIDs)
}


