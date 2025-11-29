package venues

import (
	"context"

	"github.com/google/uuid"

	"github.com/ignisrex/tix/core/internal/database"
	"github.com/ignisrex/tix/core/mappers"
	"github.com/ignisrex/tix/core/types"
)

type Repo struct {
	queries *database.Queries
}

func NewRepo(queries *database.Queries) *Repo {
	return &Repo{
		queries: queries,
	}
}

func (r *Repo) GetVenues(ctx context.Context) ([]types.Venue, error) {
	dbVenues, err := r.queries.GetVenues(ctx, database.GetVenuesParams{
		Limit:  50, // TODO: make configurable
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}
	return mappers.ToVenues(dbVenues), nil
}

func (r *Repo) CreateVenue(ctx context.Context, venue types.CreateVenueRequest) (types.Venue, error) {
	dbVenue, err := r.queries.CreateVenue(ctx, database.CreateVenueParams{
		Name:     venue.Name,
		Location: venue.Location,
		SeatMap:  venue.SeatMap,
	})
	if err != nil {
		return types.Venue{}, err
	}
	return mappers.ToVenue(dbVenue), nil
}

func (r *Repo) GetVenue(ctx context.Context, id uuid.UUID) (types.Venue, error) {
	dbVenue, err := r.queries.GetVenue(ctx, id)
	if err != nil {
		return types.Venue{}, err
	}
	return mappers.ToVenue(dbVenue), nil
}

func (r *Repo) UpdateVenue(ctx context.Context, id uuid.UUID, venue types.UpdateVenueRequest) (types.Venue, error) {
	dbVenue, err := r.queries.UpdateVenue(ctx, database.UpdateVenueParams{
		ID:       id,
		Name:     venue.Name,
		Location: venue.Location,
		SeatMap:  venue.SeatMap,
	})
	if err != nil {
		return types.Venue{}, err
	}
	return mappers.ToVenue(dbVenue), nil
}

func (r *Repo) DeleteVenue(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteVenue(ctx, id)
}

