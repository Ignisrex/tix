package main

import (
	"database/sql"
	"log"

	"github.com/ignisrex/tix/search/cmd/api"
	"github.com/ignisrex/tix/search/internal/config"
	"github.com/ignisrex/tix/search/internal/elasticsearch"
	_ "github.com/lib/pq"
)

func main() {
	addr := api.AddrFromConfig()

	dbURL := config.Envs.DBURL()
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database: ", err)
	}

	// Initialize Elasticsearch client
	var esClient *elasticsearch.Client
	esAddresses := config.Envs.ESAddresses()
	log.Printf("Attempting to connect to Elasticsearch at: %v", esAddresses)
	esClient, err = elasticsearch.NewClient(esAddresses)
	if err != nil {
		log.Printf("Warning: Failed to connect to Elasticsearch: %v. Search functionality will be limited.", err)
		esClient = nil
	} else {
		log.Printf("Successfully connected to Elasticsearch")
	}

	server := api.NewAPIServer(addr, db, esClient)
	if err := server.Run(); err != nil {
		log.Fatal("search service failed: ", err)
	}
}