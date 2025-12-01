package seed

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/ignisrex/tix/core/internal/database"
	"github.com/ignisrex/tix/core/internal/elasticsearch"
	"github.com/ignisrex/tix/core/service/events"
	"github.com/ignisrex/tix/core/service/tickets"
	"github.com/ignisrex/tix/core/service/venues"
	"github.com/ignisrex/tix/core/types"
)

type TicketAllocation struct {
	VIP      int `json:"vip"`
	GA       int `json:"ga"`
	FrontRow int `json:"front_row"`
}

type Event struct {
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	StartDate       string          `json:"start_date"`
	VenueName       string          `json:"venue_name"`
	TicketAllocation TicketAllocation `json:"ticket_allocation"`
}

type Venue struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type Data struct {
	Venues []Venue `json:"venues"`
	Events []Event `json:"events"`
}

// Run seeds venues and events+tickets using the provided DB connection and JSON file path.
func Run(ctx context.Context, db *sql.DB, esClient *elasticsearch.Client, path string) error {
	queries := database.New(db)

	venueRepo := venues.NewRepo(queries)
	venueSvc := venues.NewService(venueRepo)

	ticketRepo := tickets.NewRepo(queries)
	ticketSvc := tickets.NewService(ticketRepo)

	eventRepo := events.NewRepo(queries, db)
	eventSvc := events.NewService(eventRepo, ticketSvc, venueSvc, esClient, nil)

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var seed Data
	if err := json.NewDecoder(f).Decode(&seed); err != nil {
		return err
	}

	// Seed venues (create all from seed.json)
	for _, v := range seed.Venues {
		req := types.CreateVenueRequest{
			Name:     v.Name,
			Location: v.Location,
		}
		_, err := venueSvc.CreateVenue(ctx, req)
		if err != nil {
			log.Printf("seed: warning: failed to create venue %q: %v", v.Name, err)
		} else {
			log.Printf("seed: created venue %s", v.Name)
		}
	}

	// Get venue IDs by fetching all venues
	venueIDByName := make(map[string]uuid.UUID)
	allVenues, err := venueSvc.GetVenues(ctx)
	if err == nil {
		for _, v := range allVenues {
			venueIDByName[v.Name] = v.ID
		}
	}

	// Seed events + tickets
	for _, e := range seed.Events {
		venueID, ok := venueIDByName[e.VenueName]
		if !ok {
			log.Printf("seed: skipping event %q: venue %q not found", e.Title, e.VenueName)
			continue
		}

		start, err := time.Parse(time.RFC3339, e.StartDate)
		if err != nil {
			log.Printf("seed: skipping event %q: invalid start_date %q: %v", e.Title, e.StartDate, err)
			continue
		}

		createReq := types.CreateEventRequest{
			Title:       e.Title,
			Description: e.Description,
			StartDate:   start,
			VenueID:     venueID,
			TicketAllocation: types.TicketAllocation{
				VIP:      e.TicketAllocation.VIP,
				GA:       e.TicketAllocation.GA,
				FrontRow: e.TicketAllocation.FrontRow,
			},
		}

		ev, err := eventSvc.CreateEvent(ctx, createReq)
		if err != nil {
			log.Printf("seed: failed to create event %q: %v", e.Title, err)
			continue
		}
		log.Printf("seed: created event %s (%s) at venue %s", ev.Title, ev.ID, e.VenueName)
	}

	return nil
}


