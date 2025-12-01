package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	bookingclient "github.com/ignisrex/tix/core/internal/booking"
	"github.com/ignisrex/tix/core/internal/database"
	"github.com/ignisrex/tix/core/internal/elasticsearch"
	"github.com/ignisrex/tix/core/internal/search"
	"github.com/ignisrex/tix/core/service/booking"
	"github.com/ignisrex/tix/core/service/events"
	"github.com/ignisrex/tix/core/service/venues"
)

type APIServer struct {
	addr  string
	sqlDB *sql.DB
	q    *database.Queries
	esClient *elasticsearch.Client
	searchClient *search.Client
	bookingClient *bookingclient.Client
}

func NewAPIServer(addr string, sqlDB *sql.DB, esClient *elasticsearch.Client, searchClient *search.Client, bookingClient *bookingclient.Client) *APIServer {
	queries := database.New(sqlDB)
	return &APIServer{
		addr:  addr,
		sqlDB: sqlDB,
		q:    queries,
		esClient: esClient,
		searchClient: searchClient,
		bookingClient: bookingClient,
	}
}

func (s *APIServer) Run() error {
	
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("tix api!"))
	})

	//routes
	v1 := chi.NewRouter()
	v1.Get("/healthz", nil)

	eventHandler := events.NewHandler(s.q, s.sqlDB, s.esClient, s.searchClient, s.bookingClient)
	eventHandler.RegisterRoutes(v1)

	venueHandler := venues.NewHandler(s.q)
	venueHandler.RegisterRoutes(v1)

	bookingHandler := booking.NewHandler(s.bookingClient)
	bookingHandler.RegisterRoutes(v1)

	router.Mount("/api/v1", v1)
	return http.ListenAndServe(s.addr, router)
}