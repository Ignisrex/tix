package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/ignisrex/tix/core/cmd/api"
	bookingclient "github.com/ignisrex/tix/core/internal/booking"
	"github.com/ignisrex/tix/core/internal/config"
	"github.com/ignisrex/tix/core/internal/elasticsearch"
	"github.com/ignisrex/tix/core/internal/search"
)

func main() {

	port := config.Envs.Port
	if port == "" {
		log.Fatal("PORT not found in the env")
	}

	dbURL := config.Envs.DBURL()
	if dbURL == "" {
		log.Fatal("DB_URL not found in the env")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connection to database -> ", err)
	}
	defer conn.Close()

	// Test the connection
	if err := conn.Ping(); err != nil {
		log.Fatal("Error pinging database -> ", err)
	}

	var esClient *elasticsearch.Client
	esAddresses := config.Envs.ESAddresses()
	log.Printf("Attempting to connect to Elasticsearch at: %v", esAddresses)
	esClient, err = elasticsearch.NewClient(esAddresses)
	if err != nil {
		log.Printf("Warning: Failed to connect to Elasticsearch: %v. Continuing without search indexing.", err)
		esClient = nil
	} else {
		log.Printf("Successfully connected to Elasticsearch")
	}

	// Initialize search service client
	searchClient := search.NewClient(config.Envs.SearchServiceURL)
	log.Printf("Search service client initialized with URL: %s", config.Envs.SearchServiceURL)

	// Initialize booking service client
	bookingClient := bookingclient.NewClient(config.Envs.BookingServiceURL)
	log.Printf("Booking service client initialized with URL: %s", config.Envs.BookingServiceURL)

	server := api.NewAPIServer(":"+port, conn, esClient, searchClient, bookingClient)
	err = server.Run()
	if err != nil {
		log.Fatal("Error starting API server -> ", err)
	}

}