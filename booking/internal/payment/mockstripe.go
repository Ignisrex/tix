package payment

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
)

// TicketInfo contains ticket information needed for payment processing
type TicketInfo struct {
	ID          uuid.UUID
	PriceCents  int32
}

// ProcessPayment simulates a payment processing with mockStripe
// Takes ticket info including price and processes payment
// Returns success 90% of the time, failure 10% of the time
func ProcessPayment(ctx context.Context, ticket TicketInfo) error {
	// Simulate payment processing delay
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(50 * time.Millisecond):
	}

	// 90% success rate: if random value < 0.1 (10%), it fails
	// Using math/rand/v2 which doesn't require seeding
	if rand.Float32() < 0.1 {
		return fmt.Errorf("payment processing failed: insufficient funds for amount %d cents", ticket.PriceCents)
	}

	return nil
}

