package booking

import (
	"errors"
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
		r.Post("/locks/check", h.handleCheckLocks)
	})
}

func (h *Handler) handleReserve(w http.ResponseWriter, r *http.Request) {
	var req types.ReserveRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if len(req.TicketIDs) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ticket_ids cannot be empty"))
		return
	}

	reservedIDs, err := h.service.ReserveTickets(r.Context(), req.TicketIDs)
	if err != nil {
		status := http.StatusInternalServerError
		message := "failed to reserve tickets"

		switch {
		case errors.Is(err, ErrTicketNotFound):
			status = http.StatusNotFound
			message = "one or more tickets not found"
		case errors.Is(err, ErrTicketSold):
			status = http.StatusGone
			message = err.Error()
		case errors.Is(err, ErrTicketReserved):
			status = http.StatusConflict
			message = err.Error()
		}

		response := types.ReserveResponse{
			Success:   false,
			Message:   message,
			TicketIDs: []uuid.UUID{},
		}
		utils.WriteJSON(w, status, response)
		return
	}

	resp := types.ReserveResponse{
		Success:   true,
		Message:   "tickets reserved successfully",
		TicketIDs: reservedIDs,
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) handlePurchase(w http.ResponseWriter, r *http.Request) {
	var req types.PurchaseRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if len(req.TicketIDs) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ticket_ids cannot be empty"))
		return
	}

	purchaseID, totalCents, err := h.service.PurchaseTickets(r.Context(), req.TicketIDs)
	if err != nil {
		status := http.StatusInternalServerError
		message := "failed to purchase tickets"

		switch {
		case errors.Is(err, ErrTicketNotFound):
			status = http.StatusNotFound
			message = "one or more tickets not found"
		case errors.Is(err, ErrTicketReserved):
			status = http.StatusConflict
			message = "one or more tickets are not reserved"
		case errors.Is(err, ErrPaymentFailed):
			status = http.StatusPaymentRequired
			message = err.Error()
		}

		response := types.PurchaseResponse{
			Success:    false,
			Message:    message,
			TicketIDs:  []uuid.UUID{},
			Total:      0,
			PurchaseID: uuid.Nil,
		}
		utils.WriteJSON(w, status, response)
		return
	}

	response := types.PurchaseResponse{
		Success:    true,
		Message:    "purchase completed successfully",
		TicketIDs:  req.TicketIDs,
		Total:      totalCents,
		PurchaseID: purchaseID,
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleGetPurchase(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid purchase id: %w", err))
		return
	}

	response, err := h.service.GetPurchaseDetails(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrPurchaseNotFound) {
			status = http.StatusNotFound
		}
		utils.WriteError(w, status, fmt.Errorf("failed to get purchase details: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleCheckLocks(w http.ResponseWriter, r *http.Request) {
	var req types.CheckLocksRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if len(req.TicketIDs) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ticket_ids cannot be empty"))
		return
	}

	locks, err := h.service.CheckTicketLocks(r.Context(), req.TicketIDs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to check locks: %w", err))
		return
	}

	// Convert UUID map to string map for JSON
	locksMap := make(map[string]bool)
	for ticketID, isReserved := range locks {
		locksMap[ticketID.String()] = isReserved
	}

	response := types.CheckLocksResponse{
		Locks: locksMap,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
