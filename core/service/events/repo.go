package events

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/ignisrex/tix/core/internal/database"
	"github.com/ignisrex/tix/core/mappers"
	"github.com/ignisrex/tix/core/types"
)

type Repo struct {
	queries *database.Queries
	db *sql.DB
}

func NewRepo(queries *database.Queries, db *sql.DB) *Repo {
	return &Repo{
		queries: queries,
		db: db,
	}
}

func (r *Repo) GetEvents(ctx context.Context) ([]types.Event, error) {	
	dbEvents, err := r.queries.GetEvents(ctx, database.GetEventsParams{
		Limit:  10, //make this a configurable parameter
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}
	return mappers.ToEvents(dbEvents), nil
}

func (r *Repo) CreateEvent(ctx context.Context, event types.CreateEventRequest, tx *sql.Tx) (types.Event, error) {
	var queries *database.Queries
	if tx != nil {
		queries = r.queries.WithTx(tx)
	} else {
		queries = r.queries
	}

	dbEvent, err := queries.CreateEvent(ctx, database.CreateEventParams{
		Title: event.Title,
		Description: event.Description,
		StartDate: event.StartDate,
		VenueID: event.VenueID,
	})
	if err != nil {
		return types.Event{}, err
	}
	return mappers.ToEvent(dbEvent), nil
}

func (r *Repo) GetEvent(ctx context.Context, id uuid.UUID) (types.Event, error) {
	dbEvent, err := r.queries.GetEvent(ctx, id)
	if err != nil {
		return types.Event{}, err
	}
	return mappers.ToEvent(dbEvent), nil
}

func (r *Repo) UpdateEvent(ctx context.Context, id uuid.UUID, event types.UpdateEventRequest) (types.Event, error) {
	dbEvent, err := r.queries.UpdateEvent(ctx, database.UpdateEventParams{
		ID: id,
		Title: event.Title,
		Description: event.Description,
		StartDate: event.StartDate,
		VenueID: event.VenueID,
	})
	if err != nil {
		return types.Event{}, err
	}
	return mappers.ToEvent(dbEvent), nil
}

func (r *Repo) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteEvent(ctx, id)
	if err != nil {
		return err
	}
	return nil
}