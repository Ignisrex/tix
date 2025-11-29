package tickets

import (
	"github.com/go-chi/chi/v5"
	"github.com/ignisrex/tix/core/internal/database"
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
	r.Route("/tickets", func(r chi.Router) {
		// Add your routes here
		// r.Get("/", h.ListTickets)
		// r.Post("/", h.CreateTicket)
		// r.Get("/{id}", h.GetTicket)
		// r.Put("/{id}", h.UpdateTicket)
		// r.Delete("/{id}", h.DeleteTicket)
	})
}

// Add your HTTP handler methods here
// Example:
// func (h *Handler) CreateTicket(w http.ResponseWriter, r *http.Request) {
//     // Parse request
//     // Call service
//     // Format response
// }

