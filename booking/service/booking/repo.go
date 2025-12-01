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

