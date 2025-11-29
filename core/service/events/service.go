package events

import (
	"context"

	"github.com/google/uuid"
	
	"github.com/ignisrex/tix/core/types"
	"github.com/ignisrex/tix/core/service/tickets"
)

type Service struct {
	repo *Repo
	ticketService *tickets.Service
}

func NewService(repo *Repo, ticketService *tickets.Service) *Service {
	return &Service{
		repo: repo,
		ticketService: ticketService,
	}
}

func (s *Service) GetEvents(ctx context.Context) ([]types.Event, error) {
	return s.repo.GetEvents(ctx)
}

func (s *Service) CreateEvent(ctx context.Context, createEventRequest types.CreateEventRequest) (types.Event, error) {
	tx, err := s.repo.db.BeginTx(ctx, nil)
	if err != nil {
		return types.Event{}, err
	}
	defer tx.Rollback()

	event, err := s.repo.CreateEvent(ctx, createEventRequest, tx)
	if err != nil {
		return types.Event{}, err
	}

	_, err = s.ticketService.CreateTicketsForEvent(ctx, event.ID, createEventRequest.TicketAllocation, tx)
	if err != nil {
		return types.Event{}, err
	}

	if err := tx.Commit(); err != nil {
		return types.Event{}, err
	}

	return event, nil
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