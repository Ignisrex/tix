package events

import (
	"context"

	"github.com/ignisrex/tix/search/internal/elasticsearch"
)

type Repo struct {
	esClient *elasticsearch.Client
}

func NewRepo(esClient *elasticsearch.Client) *Repo {
	return &Repo{esClient: esClient}
}

func (r *Repo) SearchEvents(ctx context.Context, query string, limit, offset int) (*elasticsearch.SearchResponse, error) {
	if r.esClient == nil {
		return &elasticsearch.SearchResponse{Results: []elasticsearch.SearchResult{}, Total: 0}, nil
	}
	return r.esClient.SearchEvents(ctx, query, limit, offset)
}