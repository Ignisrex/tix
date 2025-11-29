package venues

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
	r.Route("/venues", func(r chi.Router) {
		r.Get("/", h.GetVenues)
		r.Post("/", h.CreateVenue)
		r.Get("/{id}", h.GetVenue)
		r.Put("/{id}", h.UpdateVenue)
		r.Delete("/{id}", h.DeleteVenue)
	})
}

func (h *Handler) GetVenues(w http.ResponseWriter, r *http.Request) {
	venues, err := h.service.GetVenues(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get venues: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, venues)
}

func (h *Handler) CreateVenue(w http.ResponseWriter, r *http.Request) {
	var req types.CreateVenueRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse create venue request body: %w", err))
		return
	}

	venue, err := h.service.CreateVenue(r.Context(), req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to create venue: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusCreated, venue)
}

func (h *Handler) GetVenue(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	venue, err := h.service.GetVenue(r.Context(), uuid.MustParse(id))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get venue: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, venue)
}

func (h *Handler) UpdateVenue(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req types.UpdateVenueRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse update venue request body: %w", err))
		return
	}

	venue, err := h.service.UpdateVenue(r.Context(), uuid.MustParse(id), req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update venue: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, venue)
}

func (h *Handler) DeleteVenue(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.DeleteVenue(r.Context(), uuid.MustParse(id)); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete venue: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("venue deleted successfully with id: %v", id))
}

