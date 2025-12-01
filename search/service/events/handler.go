package events

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ignisrex/tix/search/internal/elasticsearch"

	"github.com/ignisrex/tix/search/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(esClient *elasticsearch.Client) *Handler {
	repo := NewRepo(esClient)
	service := NewService(repo)
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/search", func(r chi.Router) {
		r.Get("/events", h.SearchEvents)
	})
}

func (h *Handler) SearchEvents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("query parameter 'q' is required"))
		return
	}

	// Parse pagination parameters
	limit := 10
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	results, err := h.service.SearchEvents(r.Context(), query, limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to search events: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, results)
}