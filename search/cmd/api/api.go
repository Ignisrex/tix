package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ignisrex/tix/search/internal/elasticsearch"

	"github.com/ignisrex/tix/search/internal/config"
	"github.com/ignisrex/tix/search/service/events"
)

type APIServer struct {
	addr    string
	db *sql.DB
	esClient *elasticsearch.Client
}

func NewAPIServer(addr string, db *sql.DB, esClient *elasticsearch.Client) *APIServer {
	return &APIServer{
		addr:    addr,
		db: db,
		esClient: esClient,
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
		w.Write([]byte("search service"))
	})

	// search endpoints under /api/v1
	v1 := chi.NewRouter()
	v1.Get("/healthz", nil)
	
	eventsHandler := events.NewHandler(s.esClient)
	eventsHandler.RegisterRoutes(v1)
	
	r.Mount("/api/v1", v1)

	return http.ListenAndServe(s.addr, r)
}

func AddrFromConfig() string {
	port := config.Envs.Port
	if port == "" {
		port = "8081"
	}
	return ":" + port
}