package booking

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ignisrex/tix/booking/internal/database"
	"github.com/ignisrex/tix/booking/internal/redis"
	"github.com/ignisrex/tix/booking/internal/utils"
	"github.com/ignisrex/tix/booking/types"
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
		r.Post("/reserve", h.handleReserve)
		r.Post("/purchase", h.handlePurchase)
		r.Get("/purchases/{id}", h.handleGetPurchase)
	})
}

func (h *Handler) handleReserve(w http.ResponseWriter, r *http.Request) {
	var req types.ReserveRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if len(req.TicketIDs) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ticket_ids cannot be empty"))
		return
	}

	response, statusCode, err := h.service.ReserveTickets(r.Context(), req.TicketIDs)
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
	var req types.PurchaseRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if len(req.TicketIDs) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ticket_ids cannot be empty"))
		return
	}

	response, statusCode, err := h.service.PurchaseTickets(r.Context(), req.TicketIDs)
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

func (h *Handler) handleGetPurchase(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid purchase id: %w", err))
		return
	}

	response, statusCode, err := h.service.GetPurchaseDetails(r.Context(), id)
	if err != nil {
		utils.WriteError(w, statusCode, fmt.Errorf("failed to get purchase details: %w", err))
		return
	}

	utils.WriteJSON(w, statusCode, response)
}
