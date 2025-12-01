package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/ignisrex/tix/core/cmd/api"
	bookingclient "github.com/ignisrex/tix/core/internal/booking"
	"github.com/ignisrex/tix/core/internal/config"
	"github.com/ignisrex/tix/core/internal/elasticsearch"
	"github.com/ignisrex/tix/core/internal/search"
	"github.com/ignisrex/tix/core/internal/seed"
)

func main() {
	ctx := context.Background()

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

	if err := conn.PingContext(ctx); err != nil {
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

	if os.Getenv("SEED_ON_START") == "true" {
		seedDatabase(ctx, conn, esClient)
	}

	
	searchClient := search.NewClient(config.Envs.SearchServiceURL)
	log.Printf("Search service client initialized with URL: %s", config.Envs.SearchServiceURL)

	
	bookingClient := bookingclient.NewClient(config.Envs.BookingServiceURL)
	log.Printf("Booking service client initialized with URL: %s", config.Envs.BookingServiceURL)

	server := api.NewAPIServer(":"+port, conn, esClient, searchClient, bookingClient)
	err = server.Run()
	if err != nil {
		log.Fatal("Error starting API server -> ", err)
	}

}

func seedDatabase(ctx context.Context, conn *sql.DB, esClient *elasticsearch.Client) {
	log.Println("SEED_ON_START=true detected, running database seeder...")
	if err := seed.Run(ctx, conn, esClient,"seed.json"); err != nil {
		log.Printf("Warning: seeding failed: %v", err)
	} else {
		log.Println("Seeding completed successfully")
	}
}