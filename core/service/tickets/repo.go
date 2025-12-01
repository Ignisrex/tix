package tickets

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/ignisrex/tix/core/internal/database"
	"github.com/ignisrex/tix/core/mappers"
	"github.com/ignisrex/tix/core/types"
)

type Repo struct {
	queries *database.Queries
	ticketTypes map[string]uuid.UUID
}

func NewRepo(queries *database.Queries) *Repo {
	return &Repo{
		queries: queries,
	}
}

func (r *Repo) GetTicket(ctx context.Context, eventID uuid.UUID, ticketID uuid.UUID) (types.Ticket, error) {
	dbTicket, err := r.queries.GetTicket(ctx, database.GetTicketParams{
		EventID: eventID,
		ID:      ticketID,
	})
	if err != nil {
		return types.Ticket{}, err
	}
	return mappers.ToEnrichedTicket(dbTicket), nil
}

func (r *Repo) GetTicketsForEvent(ctx context.Context, eventID uuid.UUID) ([]types.Ticket, error) {
	dbTickets, err := r.queries.GetTicketsForEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return mappers.ToEnrichedTickets(dbTickets), nil
}

func (r *Repo) CreateTicketsForEvent(ctx context.Context, eventID uuid.UUID, ticketTypeIDs []uuid.UUID, tx *sql.Tx) ([]types.Ticket, error) {
	
	var queries *database.Queries
	if tx != nil {
		queries = r.queries.WithTx(tx)
	} else {
		queries = r.queries
	}
	
	dbTickets, err := queries.BatchCreateTickets(ctx,
		database.BatchCreateTicketsParams{
			Column1: eventID,
			Column2: ticketTypeIDs,
		})
	if err != nil {
		return nil, err
	}
	return mappers.ToTickets(dbTickets), nil
}


func (r *Repo) getTicketTypeID(ctx context.Context, ticketTypeName string) (uuid.UUID, error) {
	if r.ticketTypes == nil {
		dbTicketTypes, err := r.queries.GetTicketTypes(ctx)
		if err != nil {
			return uuid.Nil, err
		}
		r.ticketTypes = make(map[string]uuid.UUID)
		for _, tt := range dbTicketTypes {
			r.ticketTypes[tt.Name] = tt.ID
		}
	}

	id, ok := r.ticketTypes[ticketTypeName]
	if !ok {
		return uuid.Nil, fmt.Errorf("ticket type %s not found", ticketTypeName)
	}
	return id, nil
}