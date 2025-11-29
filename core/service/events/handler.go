package events

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ignisrex/tix/core/internal/database"
	"github.com/ignisrex/tix/core/internal/utils"
	"github.com/ignisrex/tix/core/types"
)

type Handler struct {
	service *Service
}

func NewHandler(queries *database.Queries) *Handler {
	repo := NewRepo(queries)
	service := NewService(repo)
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/events", func(r chi.Router) {
		r.Get("/", h.GetEvents)
		r.Post("/", h.CreateEvent)
		r.Get("/{id}", h.GetEvent)
		r.Put("/{id}", h.UpdateEvent)
		r.Delete("/{id}", h.DeleteEvent)
	})
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetEvents(r.Context())
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

	event, err := h.service.CreateEvent(r.Context(), createEventRequest)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to create event: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusCreated, fmt.Sprintf("event created successfully with id: %v", event))
}

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	event, err := h.service.GetEvent(r.Context(), uuid.MustParse(id))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get event: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, event)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var updateEventRequest types.UpdateEventRequest
	if err := utils.ParseJSON(r, &updateEventRequest); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse update event request body: %w", err))
		return
	}
	event, err := h.service.UpdateEvent(r.Context(), uuid.MustParse(id), updateEventRequest)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update event: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, event)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.service.DeleteEvent(r.Context(), uuid.MustParse(id))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete event: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("event deleted successfully with id: %v", id))
}