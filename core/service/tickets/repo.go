package tickets

import (
	"github.com/ignisrex/tix/core/internal/database"
)

type Repo struct {
	queries *database.Queries
}

func NewRepo(queries *database.Queries) *Repo {
	return &Repo{
		queries: queries,
	}
}

// Add your data access methods here
// Example:
// func (r *Repo) CreateTicket(ctx context.Context, params database.CreateTicketParams) (database.Ticket, error) {
//     return r.queries.CreateTicket(ctx, params)
// }

