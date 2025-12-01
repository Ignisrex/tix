package events

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"

	"github.com/ignisrex/tix/core/internal/elasticsearch"
	"github.com/ignisrex/tix/core/internal/search"
	"github.com/ignisrex/tix/core/service/tickets"
	"github.com/ignisrex/tix/core/service/venues"
	"github.com/ignisrex/tix/core/types"
)

type Service struct {
	repo          *Repo
	ticketService *tickets.Service
	venueService  *venues.Service
	esClient      *elasticsearch.Client
	searchClient  *search.Client
}

func NewService(repo *Repo, ticketService *tickets.Service, venueService *venues.Service, esClient *elasticsearch.Client, searchClient *search.Client) *Service {
	return &Service{
		repo:          repo,
		ticketService: ticketService,
		venueService:  venueService,
		esClient:      esClient,
		searchClient:  searchClient,
	}
}

func (s *Service) GetEvents(ctx context.Context) ([]types.Event, error) {
	return s.repo.GetEvents(ctx)
}

func (s *Service) GetEventsWithQuery(ctx context.Context, query string, limit, offset int) ([]types.SearchEventResult, error) {
	if s.searchClient != nil {
		return s.searchClient.SearchEvents(ctx, query, limit, offset)
	}
	return nil, errors.New("search client is not available")
}

func (s *Service) CreateEvent(ctx context.Context, createEventRequest types.CreateEventRequest) (types.Event, error) {
	tx, err := s.repo.db.BeginTx(ctx, nil)
	if err != nil {
		return types.Event{}, err
	}
	defer tx.Rollback()

	event, err := s.repo.CreateEvent(ctx, createEventRequest, tx)
	if err != nil {
		log.Printf("Warning: failed to create event: %v", err)
		return types.Event{}, err
	}

	_, err = s.ticketService.CreateTicketsForEvent(ctx, event.ID, createEventRequest.TicketAllocation, tx)
	if err != nil {
		log.Printf("Warning: failed to create tickets for event: %v", err)
		return types.Event{}, err
	}

	if err := tx.Commit(); err != nil {
		return types.Event{}, err
	}


	if s.esClient != nil {
		log.Printf("Indexing event %s in Elasticsearch", event.ID)
		venue, err := s.venueService.GetVenue(ctx, event.VenueID)
		if err != nil {
			log.Printf("Warning: failed to fetch venue for ES indexing: %v", err)
			return types.Event{}, err
		}
		if err := s.esClient.IndexEvent(ctx, event, venue); err != nil {
			log.Printf("Warning: failed to index event in Elasticsearch: %v", err)
		} else {
			log.Printf("Successfully indexed event %s in Elasticsearch", event.ID)
		}
	} else {
		log.Printf("Elasticsearch client is nil, skipping indexing")
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