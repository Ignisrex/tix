package booking

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ignisrex/tix/booking/internal/database"
	"github.com/ignisrex/tix/booking/internal/redis"
	"github.com/ignisrex/tix/booking/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(queries *database.Queries, redisClient *redis.Client) *Handler {
	repo := NewRepo(queries)
	service := NewService(repo, redisClient)
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

	response, statusCode, err := h.service.Reserve(r.Context(), id)
	if err != nil {
		// If response exists but success is false, it's an "already reserved" error
		if response != nil && !response.Success {
			utils.WriteJSON(w, statusCode, response)
			return
		}
		utils.WriteError(w, statusCode, fmt.Errorf("failed to reserve: %w", err))
		return
	}

	utils.WriteJSON(w, statusCode, response)
}

func (h *Handler) handlePurchase(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id: %w", err))
		return
	}

	response, statusCode, err := h.service.Purchase(r.Context(), id)
	if err != nil {
		// If response exists but success is false, it's a payment failure or not found
		if response != nil && !response.Success {
			utils.WriteJSON(w, statusCode, response)
			return
		}
		utils.WriteError(w, statusCode, fmt.Errorf("failed to purchase: %w", err))
		return
	}

	utils.WriteJSON(w, statusCode, response)
}
