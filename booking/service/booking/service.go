package booking

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Reserve(ctx context.Context, id uuid.UUID) error {
	return s.repo.ReserveTicket(ctx, id)
}

func (s *Service) Purchase(ctx context.Context, id uuid.UUID) error {
	return s.repo.PurchaseTicket(ctx, id)
}


