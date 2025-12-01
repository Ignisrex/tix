package booking

import (
	"context"

	"github.com/google/uuid"

	"github.com/ignisrex/tix/booking/internal/database"
)

type Repo struct {
	queries *database.Queries
}

func NewRepo(queries *database.Queries) *Repo {
	return &Repo{queries: queries}
}

func (r *Repo) GetTicketStatus(ctx context.Context, ticketID uuid.UUID) (database.TicketStatus, error) {
	return r.queries.GetTicketStatus(ctx, ticketID)
}

func (r *Repo) GetTicketWithPrice(ctx context.Context, ticketID uuid.UUID) (database.GetTicketWithPriceRow, error) {
	return r.queries.GetTicketWithPrice(ctx, ticketID)
}

func (r *Repo) PurchaseTicket(ctx context.Context, ticketID uuid.UUID) error {
	return r.queries.PurchaseTicket(ctx, ticketID)
}

func (r *Repo) PurchaseTickets(ctx context.Context, ticketIDs []uuid.UUID, totalCents int32) (uuid.UUID, error) {
	return r.queries.PurchaseTickets(ctx, database.PurchaseTicketsParams{
		TotalCents: totalCents,
		Column2:    ticketIDs,
	})
}

func (r *Repo) GetPurchaseDetails(ctx context.Context, purchaseID uuid.UUID) (database.GetPurchaseDetailsRow, error) {
	return r.queries.GetPurchaseDetails(ctx, purchaseID)
}

