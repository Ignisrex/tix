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

func (r *Repo) ReserveTicket(ctx context.Context, ticketID uuid.UUID) error {
	return r.queries.ReserveTicket(ctx, ticketID)
}

func (r *Repo) PurchaseTicket(ctx context.Context, ticketID uuid.UUID) error {
	return r.queries.PurchaseTicket(ctx, ticketID)
}

