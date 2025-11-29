package main

import (
	"database/sql"
	"log"

	"github.com/ignisrex/tix/booking/cmd/api"
	"github.com/ignisrex/tix/booking/internal/config"
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

	server := api.NewAPIServer(addr, db)
	if err := server.Run(); err != nil {
		log.Fatal("booking service failed: ", err)
	}
}