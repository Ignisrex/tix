package types

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TicketStatus string
const (
	TicketStatusAvailable TicketStatus = "available"
	TicketStatusSold      TicketStatus = "sold"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	StartDate   time.Time `json:"start_date"`
	VenueID     uuid.UUID `json:"venue_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
}

type Venue struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name" validate:"required"`
	Location string    `json:"location" validate:"required"`
}

type Ticket struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id" validate:"required"`
	TicketTypeID uuid.UUID `json:"ticket_type_id" validate:"required"`
	Status    TicketStatus `json:"status"`
}

type CreateEventRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	VenueID     uuid.UUID `json:"venue_id" validate:"required"`
	TicketAllocation TicketAllocation `json:"ticket_allocation" validate:"required"`
}

type CreateVenueRequest struct {
	Name     string    `json:"name" validate:"required"`
	Location string    `json:"location" validate:"required"`
}

type TicketAllocation struct {
	VIP int `json:"vip" validate:"required"`
	GA int `json:"ga" validate:"required"`
	FrontRow int `json:"front_row" validate:"required"`
}

type UpdateEventRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	VenueID     uuid.UUID `json:"venue_id" validate:"required"`
}

type UpdateVenueRequest struct {
	Name     string    `json:"name" validate:"required"`
	Location string    `json:"location" validate:"required"`
	SeatMap  json.RawMessage `json:"seat_map" validate:"required"`
}	
type SearchEventResult struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	VenueID        string    `json:"venue_id"`
	VenueName      string    `json:"venue_name"`
	VenueLocation  string    `json:"venue_location"`
	CreatedAt      time.Time `json:"created_at"`
}

type SearchEventResults struct {
	Results []SearchEventResult `json:"results"`
	Total   int                 `json:"total"`
}
