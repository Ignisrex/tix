package events

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

func (s *Service) GetEvents(ctx context.Context) ([]types.Event, error) {
	return s.repo.GetEvents(ctx)
}

func (s *Service) CreateEvent(ctx context.Context, event types.CreateEventRequest) (types.Event, error) {
	return s.repo.CreateEvent(ctx, event)
}

func (s *Service) GetEvent(ctx context.Context, id uuid.UUID) (types.Event, error) {
	return s.repo.GetEvent(ctx, id)
}

func (s *Service) UpdateEvent(ctx context.Context, id uuid.UUID, event types.UpdateEventRequest) (types.Event, error) {
	return s.repo.UpdateEvent(ctx, id, event)
}

func (s *Service) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteEvent(ctx, id)
}