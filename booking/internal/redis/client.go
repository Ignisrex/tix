package redis

import (
	"context"
	"errors"
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

var reserveScript = redis.NewScript(`
for i = 1, #KEYS do
  if redis.call("EXISTS", KEYS[i]) == 1 then
    return 0
  end
end
for i = 1, #KEYS do
  local result = redis.call("SET", KEYS[i], "true", "EX", ARGV[1], "NX")
  if result == false then
    -- SET NX failed (key was created by another process between EXISTS check and SET)
    return 0
  end
end
return 1
`)

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

/*ReserveTickets attempts to reserve multiple tickets atomically during the operation*/
func (c *Client) ReserveTickets(ctx context.Context, ticketIDs []uuid.UUID) error {
	keys := make([]string, len(ticketIDs))
	for i, id := range ticketIDs {
		keys[i] = keyPrefix + id.String()
	}

	res, err := reserveScript.Run(ctx, c.rdb, keys, int(c.ttl.Seconds())).Int()
	if err != nil {
		return fmt.Errorf("failed to reserve tickets: %w", err)
	}

	if res == 0 {
		return errors.New("failed to reserve tickets. At least one ticket is already reserved")
	}

	return nil
}

func (c *Client) RefreshTickets(ctx context.Context, ticketIDs []uuid.UUID, ttl time.Duration) (bool, error) {
	pipe := c.rdb.Pipeline()
	for _, ticketID := range ticketIDs {
		key := keyPrefix + ticketID.String()
		pipe.Expire(ctx, key, ttl)
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to refresh tickets: %w", err)
	}

	for _, cmd := range cmds {
		if cmd.Err() != nil {
			return false, fmt.Errorf("failed to refresh ticket: %w", cmd.Err())
		}

		expireCmd, ok := cmd.(*redis.BoolCmd)
		if !ok {
			return false, fmt.Errorf("failed to refresh ticket: %w", cmd.Err())
		}

		if !expireCmd.Val() {
			return false, fmt.Errorf("failed to refresh ticket: %w", expireCmd.Err())
		}
	}
	return true, nil
}

func (c *Client) ReleaseTickets(ctx context.Context, ticketIDs []uuid.UUID) error {
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


// Returns a map of ticketID -> is_reserved (true if reserved, false if available).
func (c *Client) AreReserved(ctx context.Context, ticketIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	results := make(map[uuid.UUID]bool)
	pipe := c.rdb.Pipeline()
	for _, ticketID := range ticketIDs {
		key := keyPrefix + ticketID.String()
		pipe.Exists(ctx, key)
	}

	// Execute all operations atomically
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check reservations: %w", err)
	}

	// Process results
	for i, ticketID := range ticketIDs {
		existsCmd, ok := cmds[i].(*redis.IntCmd)
		if !ok {
			return nil, fmt.Errorf("failed to check reservation: %w", cmds[i].Err())
		}
		if existsCmd.Val() > 0 {
			results[ticketID] = true
			continue
		}
		results[ticketID] = false
	}

	return results, nil
}

