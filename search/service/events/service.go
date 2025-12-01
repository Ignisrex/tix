package events

import (
	"context"

	"github.com/ignisrex/tix/search/internal/elasticsearch"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) SearchEvents(ctx context.Context, query string, limit, offset int) (*elasticsearch.SearchResponse, error) {
	return s.repo.SearchEvents(ctx, query, limit, offset)
}