package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/ignisrex/tix/core/cmd/api"
	"github.com/ignisrex/tix/core/internal/config"
	"github.com/ignisrex/tix/core/internal/database"
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

	queries := database.New(conn)

	server := api.NewAPIServer(":"+port, conn,queries)
	err = server.Run()
	if err != nil {
		log.Fatal("Error starting API server -> ", err)
	}

}