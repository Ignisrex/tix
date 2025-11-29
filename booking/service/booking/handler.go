package booking

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ignisrex/tix/booking/internal/database"
	"github.com/ignisrex/tix/booking/internal/utils"
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
	r.Route("/booking", func(r chi.Router) {
		r.Post("/reserve/{id}", h.handleReserve)
		r.Post("/purchase/{id}", h.handlePurchase)
	})
}

func (h *Handler) handleReserve(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id: %w", err))
		return
	}

	if err := h.service.Reserve(r.Context(), id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to reserve: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message":   "reservation created",
		"ticket_id": id,
	})
}

func (h *Handler) handlePurchase(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id: %w", err))
		return
	}

	if err := h.service.Purchase(r.Context(), id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to purchase: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message":   "purchase completed",
		"ticket_id": id,
	})
}
