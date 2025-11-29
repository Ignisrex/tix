package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/ignisrex/tix/booking/internal/config"
	"github.com/ignisrex/tix/booking/internal/database"
	"github.com/ignisrex/tix/booking/service/booking"
)

type APIServer struct {
	addr    string
	db *sql.DB
	queries *database.Queries
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	queries := database.New(db)
	return &APIServer{
		addr:    addr,
		db: db,
		queries: queries,
	}
}

func (s *APIServer) Run() error {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("booking service"))
	})

	// booking endpoints under /api/v1
	v1 := chi.NewRouter()
	bookingHandler := booking.NewHandler(s.queries)
	bookingHandler.RegisterRoutes(v1)
	r.Mount("/api/v1", v1)

	return http.ListenAndServe(s.addr, r)
}

// AddrFromConfig builds the listen address from env config.
func AddrFromConfig() string {
	port := config.Envs.Port
	if port == "" {
		port = "8081"
	}
	return ":" + port
}