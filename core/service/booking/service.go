package booking

import (
	"context"

	"github.com/google/uuid"

	bookingclient "github.com/ignisrex/tix/core/internal/booking"
)

type Service struct {
	bookingClient *bookingclient.Client
}

func NewService(bookingClient *bookingclient.Client) *Service {
	return &Service{
		bookingClient: bookingClient,
	}
}

func (s *Service) ReserveTickets(ctx context.Context, ticketIDs []uuid.UUID) (*bookingclient.ReserveResponse, int, error) {
	return s.bookingClient.ReserveTickets(ctx, ticketIDs)
}

func (s *Service) PurchaseTickets(ctx context.Context, ticketIDs []uuid.UUID) (*bookingclient.PurchaseResponse, int, error) {
	return s.bookingClient.PurchaseTickets(ctx, ticketIDs)
}

func (s *Service) GetPurchaseDetails(ctx context.Context, purchaseID uuid.UUID) (*bookingclient.PurchaseDetailsResponse, int, error) {
	return s.bookingClient.GetPurchaseDetails(ctx, purchaseID)
}

