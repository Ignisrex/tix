package mappers

import (
	"github.com/ignisrex/tix/core/internal/database"
	"github.com/ignisrex/tix/core/types"
)

func ToEvent(dbEvent database.Event) types.Event {
	return types.Event{
		ID: dbEvent.ID,
		Title: dbEvent.Title,
		Description: dbEvent.Description,
		StartDate: dbEvent.StartDate,
		VenueID: dbEvent.VenueID,
	}
}

func ToEvents(dbEvents []database.Event) []types.Event {
	events := make([]types.Event, len(dbEvents))
	for i, dbEvent := range dbEvents {
		events[i] = ToEvent(dbEvent)
	}
	return events
}

func ToVenue(dbVenue database.Venue) types.Venue {
	return types.Venue{
		ID: dbVenue.ID,
		Name: dbVenue.Name,
		Location: dbVenue.Location,
	}
}

func ToVenues(dbVenues []database.Venue) []types.Venue {
	venues := make([]types.Venue, len(dbVenues))
	for i, dbVenue := range dbVenues {
		venues[i] = ToVenue(dbVenue)
	}
	return venues
}

func ToTicket(dbTicket database.Ticket) types.Ticket {
	return types.Ticket{	
		ID: dbTicket.ID,
		EventID: dbTicket.EventID,
		TicketTypeID: dbTicket.TicketTypeID,
		Status: dbTicket.Status,
	}
}

func ToTickets(dbTickets []database.Ticket) []types.Ticket {
	tickets := make([]types.Ticket, len(dbTickets))
	for i, dbTicket := range dbTickets {
		tickets[i] = ToTicket(dbTicket)
	}
	return tickets
}