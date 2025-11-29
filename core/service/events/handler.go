package events

import (
	"fmt"
	"net/http"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ignisrex/tix/core/internal/database"
	"github.com/ignisrex/tix/core/internal/utils"
	"github.com/ignisrex/tix/core/types"
	"github.com/ignisrex/tix/core/service/tickets"
)

type Handler struct {
	eventService *Service
	ticketService *tickets.Service
}

func NewHandler(queries *database.Queries, db *sql.DB) *Handler {
	ticketRepo := tickets.NewRepo(queries)
	ticketService := tickets.NewService(ticketRepo)

	eventRepo := NewRepo(queries, db)
	eventService := NewService(eventRepo, ticketService)
	
	return &Handler{
		eventService: eventService,
		ticketService: ticketService,
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
		})
	})
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.eventService.GetEvents(r.Context())
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