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
		r.Post("/reserve/{ticket_id}", h.ReserveTicket)
		r.Post("/purchase/{ticket_id}", h.PurchaseTicket)
	})
}

func (h *Handler) ReserveTicket(w http.ResponseWriter, r *http.Request) {
	ticketIDStr := chi.URLParam(r, "ticket_id")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid ticket_id: %w", err))
		return
	}

	response, statusCode, err := h.service.ReserveTicket(r.Context(), ticketID)
	if err != nil {
		if response != nil && !response.Success {
			utils.WriteJSON(w, statusCode, response)
			return
		}
		utils.WriteError(w, statusCode, fmt.Errorf("failed to reserve ticket: %w", err))
		return
	}

	utils.WriteJSON(w, statusCode, response)
}

func (h *Handler) PurchaseTicket(w http.ResponseWriter, r *http.Request) {
	ticketIDStr := chi.URLParam(r, "ticket_id")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid ticket_id: %w", err))
		return
	}

	response, statusCode, err := h.service.PurchaseTicket(r.Context(), ticketID)
	if err != nil {
		if response != nil && !response.Success {
			utils.WriteJSON(w, statusCode, response)
			return
		}
		utils.WriteError(w, statusCode, fmt.Errorf("failed to purchase ticket: %w", err))
		return
	}

	utils.WriteJSON(w, statusCode, response)
}

