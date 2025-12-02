package mappers

import (
	"github.com/ignisrex/tix/booking/internal/database"
	"github.com/ignisrex/tix/booking/types"
)

func ToTicket(dbTicket database.GetTicketsWithPriceRow) types.Ticket {
	return types.Ticket{
		ID: dbTicket.ID,
	}
}

func ToTickets(dbTickets []database.GetTicketsWithPriceRow) []types.Ticket {
	tickets := make([]types.Ticket, len(dbTickets))
	for i, dbTicket := range dbTickets {
		tickets[i] = ToTicket(dbTicket)
	}
	return tickets
}