package types

import (
	"github.com/google/uuid"
)

type ReserveRequest struct {
	TicketIDs []uuid.UUID `json:"ticket_ids"`
}

type ReserveResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	TicketIDs []uuid.UUID `json:"ticket_ids"` // IDs of successfully reserved tickets
}

type PurchaseRequest struct {
	TicketIDs []uuid.UUID `json:"ticket_ids"`
}

type PurchaseResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	TicketIDs  []uuid.UUID `json:"ticket_ids"`  // IDs of successfully purchased tickets
	Total      int32       `json:"total"`       // Total price in cents
	PurchaseID uuid.UUID   `json:"purchase_id"` // ID of the purchase record
}

type PurchaseTicketDetail struct {
	ID                    uuid.UUID `json:"id"`
	EventID               uuid.UUID `json:"event_id"`
	TicketTypeID          uuid.UUID `json:"ticket_type_id"`
	Status                string    `json:"status"`
	TicketTypeName        string    `json:"ticket_type_name"`
	TicketTypeDisplayName string    `json:"ticket_type_display_name"`
	TicketTypePriceCents  int32     `json:"ticket_type_price_cents"`
}

type PurchaseDetailsResponse struct {
	PurchaseID        uuid.UUID             `json:"purchase_id"`
	TotalCents        int32                 `json:"total_cents"`
	PurchaseCreatedAt string                `json:"purchase_created_at"` // ISO timestamp
	Tickets           []PurchaseTicketDetail `json:"tickets"`
}

