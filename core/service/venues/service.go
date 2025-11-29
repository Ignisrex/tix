package venues

import (
	"context"

	"github.com/google/uuid"

	"github.com/ignisrex/tix/core/types"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetVenues(ctx context.Context) ([]types.Venue, error) {
	return s.repo.GetVenues(ctx)
}

func (s *Service) CreateVenue(ctx context.Context, venue types.CreateVenueRequest) (types.Venue, error) {
	return s.repo.CreateVenue(ctx, venue)
}

func (s *Service) GetVenue(ctx context.Context, id uuid.UUID) (types.Venue, error) {
	return s.repo.GetVenue(ctx, id)
}

func (s *Service) UpdateVenue(ctx context.Context, id uuid.UUID, venue types.UpdateVenueRequest) (types.Venue, error) {
	return s.repo.UpdateVenue(ctx, id, venue)
}

func (s *Service) DeleteVenue(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteVenue(ctx, id)
}

