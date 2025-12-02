package booking

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	bookingclient "github.com/ignisrex/tix/core/internal/booking"
	"github.com/ignisrex/tix/core/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(bookingClient *bookingclient.Client) *Handler {
	service := NewService(bookingClient)
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/booking", func(r chi.Router) {
		r.Post("/reserve", h.ReserveTickets)
		r.Post("/purchase", h.PurchaseTickets)
		r.Get("/purchases/{id}", h.GetPurchaseDetails)
	})
}

func (h *Handler) ReserveTickets(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TicketIDs []uuid.UUID `json:"ticket_ids"`
	}
	
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if len(req.TicketIDs) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ticket_ids cannot be empty"))
		return
	}

	response, statusCode, err := h.service.ReserveTickets(r.Context(), req.TicketIDs)
	if err != nil {
		if response != nil && !response.Success {
			utils.WriteJSON(w, statusCode, response)
			return
		}
		utils.WriteError(w, statusCode, fmt.Errorf("failed to reserve tickets: %w", err))
		return
	}

	utils.WriteJSON(w, statusCode, response)
}

func (h *Handler) PurchaseTickets(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TicketIDs []uuid.UUID `json:"ticket_ids"`
	}
	
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if len(req.TicketIDs) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ticket_ids cannot be empty"))
		return
	}

	response, statusCode, err := h.service.PurchaseTickets(r.Context(), req.TicketIDs)
	if err != nil {
		if response != nil && !response.Success {
			utils.WriteJSON(w, statusCode, response)
			return
		}
		utils.WriteError(w, statusCode, fmt.Errorf("failed to purchase tickets: %w", err))
		return
	}

	_ = utils.WriteJSON(w, statusCode, response)
}

func (h *Handler) GetPurchaseDetails(w http.ResponseWriter, r *http.Request) {
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

	_ = utils.WriteJSON(w, statusCode, response)
}

