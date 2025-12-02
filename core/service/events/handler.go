package events

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	bookingclient "github.com/ignisrex/tix/core/internal/booking"
	"github.com/ignisrex/tix/core/internal/database"
	"github.com/ignisrex/tix/core/internal/elasticsearch"
	"github.com/ignisrex/tix/core/internal/search"
	"github.com/ignisrex/tix/core/internal/utils"
	"github.com/ignisrex/tix/core/service/tickets"
	"github.com/ignisrex/tix/core/service/venues"
	"github.com/ignisrex/tix/core/types"
)

type Handler struct {
	eventService   *Service
	ticketService  *tickets.Service
	bookingClient  *bookingclient.Client
}

func NewHandler(queries *database.Queries, db *sql.DB, esClient *elasticsearch.Client, searchClient *search.Client, bookingClient *bookingclient.Client) *Handler {
	ticketRepo := tickets.NewRepo(queries)
	ticketService := tickets.NewService(ticketRepo)

	venueRepo := venues.NewRepo(queries)
	venueService := venues.NewService(venueRepo)

	eventRepo := NewRepo(queries, db)
	eventService := NewService(eventRepo, ticketService, venueService, esClient, searchClient)
	
	return &Handler{
		eventService:  eventService,
		ticketService: ticketService,
		bookingClient: bookingClient,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/events", func(r chi.Router) {
		r.Get("/", h.GetEvents)
		r.Post("/", h.CreateEvent)
		r.Get("/{event_id}", h.GetEvent)
		r.Put("/{event_id}", h.UpdateEvent)
		r.Delete("/{event_id}", h.DeleteEvent)

		r.Route("/{event_id}/tickets", func(r chi.Router) {
			r.Get("/", h.GetTickets)
			r.Get("/stream", h.StreamTickets)
			r.Get("/{ticket_id}", h.GetTicket)
		})
	})
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	
	query := r.URL.Query().Get("q")
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	
	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}
	
	events, err := h.eventService.GetEventsWithQuery(r.Context(), query, limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get events: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, events)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var createEventRequest types.CreateEventRequest
	if err := utils.ParseJSON(r, &createEventRequest); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse create event request body: %w", err))
		return
	}

	event, err := h.eventService.CreateEvent(r.Context(), createEventRequest)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to create event: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusCreated, fmt.Sprintf("event created successfully with id: %v", event))
}

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")
	event, err := h.eventService.GetEvent(r.Context(), uuid.MustParse(id))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get event: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, event)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")
	var updateEventRequest types.UpdateEventRequest
	if err := utils.ParseJSON(r, &updateEventRequest); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse update event request body: %w", err))
		return
	}
	event, err := h.eventService.UpdateEvent(r.Context(), uuid.MustParse(id), updateEventRequest)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update event: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, event)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")
	err := h.eventService.DeleteEvent(r.Context(), uuid.MustParse(id))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete event: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("event deleted successfully with id: %v", id))
}

func (h *Handler) GetTickets(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "event_id")
	tickets, err := h.ticketService.GetTicketsForEvent(r.Context(), uuid.MustParse(eventID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get tickets for event: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, tickets)
}

func (h *Handler) GetTicket(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "event_id")
	ticketID := chi.URLParam(r, "ticket_id")
	ticket, err := h.ticketService.GetTicket(r.Context(), uuid.MustParse(eventID), uuid.MustParse(ticketID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get ticket: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, ticket)
}

// enrichTicketsWithLocks adds lock status to tickets by checking the booking service
func (h *Handler) enrichTicketsWithLocks(ctx context.Context, tickets []types.Ticket) ([]types.Ticket, error) {
	if h.bookingClient == nil {
		// If booking client is not available, return tickets without lock status
		log.Printf("Warning: booking client is not available, returning tickets without lock status")
		return tickets, nil
	}

	if len(tickets) == 0 {
		return tickets, nil
	}

	// Collect all ticket IDs
	ticketIDs := make([]uuid.UUID, 0, len(tickets))
	for _, ticket := range tickets {
		ticketIDs = append(ticketIDs, ticket.ID)
	}

	// Check lock status for all tickets
	locks, _, err := h.bookingClient.CheckTicketLocks(ctx, ticketIDs)
	if err != nil {
		log.Printf("Warning: failed to check ticket locks: %v", err)
		// Continue without lock status if check fails
		return tickets, nil
	}

	// Enrich tickets with lock status
	enriched := make([]types.Ticket, len(tickets))
	for i, ticket := range tickets {
		enriched[i] = ticket
		enriched[i].IsReserved = locks[ticket.ID]
	}

	return enriched, nil
}


func (h *Handler) StreamTickets(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "event_id")
	eventUUID := uuid.MustParse(eventID)

	// Set up SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// Create a context that cancels when client disconnects
	ctx := r.Context()

	// Create a ticker for 2-second updates;TODO: should be configurable
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Send initial data immediately
	h.sendTicketUpdate(w, ctx, eventUUID)

	// Stream updates every 2 seconds
	for {
		select {
		case <-ctx.Done():
			// Client disconnected
			return
		case <-ticker.C:
			if err := h.sendTicketUpdate(w, ctx, eventUUID); err != nil {
				log.Printf("Error sending ticket update: %v", err)
				return
			}
			// Flush the response to ensure data is sent immediately
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}
}

func (h *Handler) sendTicketUpdate(w http.ResponseWriter, ctx context.Context, eventID uuid.UUID) error {
	// Fetch tickets from database - could overwhelm the database, if we have read replicas this would mitigate the issue
	tickets, err := h.ticketService.GetTicketsForEvent(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get tickets: %w", err)
	}

	enrichedTickets, _ := h.enrichTicketsWithLocks(ctx, tickets)
	
	jsonData, err := json.Marshal(enrichedTickets)
	if err != nil {
		return fmt.Errorf("failed to marshal tickets: %w", err)
	}

	// Write SSE formatted data
	_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
	if err != nil {
		return fmt.Errorf("failed to write SSE data: %w", err)
	}

	return nil
}