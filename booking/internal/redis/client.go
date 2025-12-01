package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	reservationTTL = 3 * time.Minute // 180 seconds
	keyPrefix      = "ticket:"
)

type Client struct {
	rdb *redis.Client
}

func NewClient(addr string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	log.Printf("Redis connection successful to %s", addr)
	return &Client{rdb: rdb}, nil
}

func (c *Client) IsReserved(ctx context.Context, ticketID uuid.UUID) (bool, error) {
	key := keyPrefix + ticketID.String()
	exists, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check reservation: %w", err)
	}
	return exists > 0, nil
}

func (c *Client) ReserveTicket(ctx context.Context, ticketID uuid.UUID) (bool, error) {
	key := keyPrefix + ticketID.String()
	
	// Use SETNX for atomic reservation (only set if not exists)
	set, err := c.rdb.SetNX(ctx, key, "true", reservationTTL).Result()
	if err != nil {
		return false, fmt.Errorf("failed to reserve ticket: %w", err)
	}
	
	return set, nil // Returns true if key was set (not already reserved), false if already exists
}

func (c *Client) ReleaseTicket(ctx context.Context, ticketID uuid.UUID) error {
	key := keyPrefix + ticketID.String()
	if err := c.rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to release ticket: %w", err)
	}
	return nil
}

