package types

import (
	"github.com/google/uuid"
)

type ReserveResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	TicketID uuid.UUID `json:"ticket_id"`
}

type PurchaseResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	TicketID uuid.UUID `json:"ticket_id"`
	Total    int32     `json:"total"` // Price in cents
}

