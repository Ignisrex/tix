package tickets

import (
	"context"
	"database/sql"

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


func (s *Service) GetTicketsForEvent(ctx context.Context, eventID uuid.UUID) ([]types.Ticket, error) {
	return s.repo.GetTicketsForEvent(ctx, eventID)
}

func (s *Service) CreateTicketsForEvent(ctx context.Context, eventID uuid.UUID, ticketAllocation types.TicketAllocation, tx *sql.Tx) ([]types.Ticket, error) {

	ticketTypeIDs := []uuid.UUID{}
	vipTicketTypeID, err := s.repo.getTicketTypeID(ctx, "vip")
	if err != nil {
		return nil, err
	}
	for i := 0; i < ticketAllocation.VIP; i++ {
		ticketTypeIDs = append(ticketTypeIDs, vipTicketTypeID)
	}

	gaTicketTypeID, err := s.repo.getTicketTypeID(ctx, "ga")
	if err != nil {
		return nil, err
	}
	for i := 0; i < ticketAllocation.GA; i++ {
		ticketTypeIDs = append(ticketTypeIDs, gaTicketTypeID)
	}

	frontRowTicketTypeID, err := s.repo.getTicketTypeID(ctx, "front_row")
	if err != nil {
		return nil, err
	}
	for i := 0; i < ticketAllocation.FrontRow; i++ {
		ticketTypeIDs = append(ticketTypeIDs, frontRowTicketTypeID)
	}

	return s.repo.CreateTicketsForEvent(ctx, eventID, ticketTypeIDs, tx)
}
