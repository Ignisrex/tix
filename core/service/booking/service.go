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

func (s *Service) ReserveTicket(ctx context.Context, ticketID uuid.UUID) (*bookingclient.ReserveResponse, int, error) {
	return s.bookingClient.ReserveTicket(ctx, ticketID)
}

func (s *Service) PurchaseTicket(ctx context.Context, ticketID uuid.UUID) (*bookingclient.PurchaseResponse, int, error) {
	return s.bookingClient.PurchaseTicket(ctx, ticketID)
}

