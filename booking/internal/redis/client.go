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
	keyPrefix = "ticket:"
)

type Client struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewClient(addr string, ttlSeconds int) (*Client, error) {
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
	return &Client{
		rdb: rdb,
		ttl: time.Duration(ttlSeconds) * time.Second,
	}, nil
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
	set, err := c.rdb.SetNX(ctx, key, "true", c.ttl).Result()
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

// ReserveTickets attempts to reserve multiple tickets atomically.
// Returns a map of ticketID -> success (true if reserved, false if already reserved)
// and any error that occurred during the operation.
func (c *Client) ReserveTickets(ctx context.Context, ticketIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	if len(ticketIDs) == 0 {
		return make(map[uuid.UUID]bool), nil
	}

	results := make(map[uuid.UUID]bool)
	pipe := c.rdb.Pipeline()

	// Prepare all SETNX operations
	ops := make(map[string]uuid.UUID)
	for _, ticketID := range ticketIDs {
		key := keyPrefix + ticketID.String()
		ops[key] = ticketID
		pipe.SetNX(ctx, key, "true", c.ttl)
	}

	// Execute all operations atomically
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to reserve tickets: %w", err)
	}

	// Process results
	i := 0
	for _, ticketID := range ticketIDs {
		if i < len(cmds) {
			setCmd, ok := cmds[i].(*redis.BoolCmd)
			if ok {
				results[ticketID] = setCmd.Val()
			} else {
				results[ticketID] = false
			}
		} else {
			results[ticketID] = false
		}
		i++
	}

	return results, nil
}

// ReleaseTickets releases multiple tickets.
func (c *Client) ReleaseTickets(ctx context.Context, ticketIDs []uuid.UUID) error {
	if len(ticketIDs) == 0 {
		return nil
	}

	pipe := c.rdb.Pipeline()
	for _, ticketID := range ticketIDs {
		key := keyPrefix + ticketID.String()
		pipe.Del(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to release tickets: %w", err)
	}
	return nil
}

