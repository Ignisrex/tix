package booking

import (
	"context"

	"github.com/google/uuid"

	"github.com/ignisrex/tix/booking/internal/database"
	"github.com/ignisrex/tix/booking/types"
	"github.com/ignisrex/tix/booking/mappers"
)

type Repo struct {
	queries *database.Queries
}

func NewRepo(queries *database.Queries) *Repo {
	return &Repo{queries: queries}
}

func (r *Repo) GetTicketsWithPrice(ctx context.Context, ticketID []uuid.UUID) ([]types.Ticket, error) {
	dbTickets, err := r.queries.GetTicketsWithPrice(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	return mappers.ToTickets(dbTickets), nil
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

