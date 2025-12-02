package main

import (
	"database/sql"
	"log"
	

	"github.com/ignisrex/tix/booking/cmd/api"
	"github.com/ignisrex/tix/booking/internal/config"
	"github.com/ignisrex/tix/booking/internal/redis"
	_ "github.com/lib/pq"
)

func main() {
	addr := api.AddrFromConfig()

	dbURL := config.Envs.DBURL()
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to create database handle: ", err)
	}
	defer db.Close()

	redisAddr := config.Envs.RedisAddr()
	redisClient, err := redis.NewClient(redisAddr, config.Envs.ReservationTTLSeconds)
	if err != nil {
		log.Fatal("failed to connect to Redis: ", err)
	}
	log.Printf("Successfully connected to Redis at %s", redisAddr)

	server := api.NewAPIServer(addr, db, redisClient)
	if err := server.Run(); err != nil {
		log.Fatal("booking service failed: ", err)
	}
}