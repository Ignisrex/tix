package booking

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ignisrex/tix/core/internal/utils"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type ReserveRequest struct {
	TicketIDs []uuid.UUID `json:"ticket_ids"`
}

type ReserveResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	TicketIDs []uuid.UUID `json:"ticket_ids"`
}

type PurchaseRequest struct {
	TicketIDs []uuid.UUID `json:"ticket_ids"`
}

type PurchaseResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	TicketIDs  []uuid.UUID `json:"ticket_ids"`
	Total      int32       `json:"total"`
	PurchaseID uuid.UUID   `json:"purchase_id"`
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
	PurchaseCreatedAt string                 `json:"purchase_created_at"`
	Tickets           []PurchaseTicketDetail `json:"tickets"`
}

func (c *Client) ReserveTickets(ctx context.Context, ticketIDs []uuid.UUID) (*ReserveResponse, int, error) {
	url := fmt.Sprintf("%s/api/v1/booking/reserve", c.baseURL)
	
	reqBody := ReserveRequest{
		TicketIDs: ticketIDs,
	}
	
	req, err := utils.MakeJSONRequest(ctx, "POST", url, reqBody)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	
	body, statusCode, err := utils.ExecuteRequest(c.httpClient, req)
	if err != nil {
		return nil, statusCode, err
	}
	
	return utils.UnmarshalJSONResponse[ReserveResponse](body, statusCode, "booking service")
}

func (c *Client) PurchaseTickets(ctx context.Context, ticketIDs []uuid.UUID) (*PurchaseResponse, int, error) {
	url := fmt.Sprintf("%s/api/v1/booking/purchase", c.baseURL)
	
	reqBody := PurchaseRequest{
		TicketIDs: ticketIDs,
	}
	
	req, err := utils.MakeJSONRequest(ctx, "POST", url, reqBody)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	
	body, statusCode, err := utils.ExecuteRequest(c.httpClient, req)
	if err != nil {
		return nil, statusCode, err
	}
	
	return utils.UnmarshalJSONResponse[PurchaseResponse](body, statusCode, "booking service")
}

func (c *Client) GetPurchaseDetails(ctx context.Context, purchaseID uuid.UUID) (*PurchaseDetailsResponse, int, error) {
	url := fmt.Sprintf("%s/api/v1/booking/purchases/%s", c.baseURL, purchaseID.String())
	
	req, err := utils.MakeJSONRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	
	body, statusCode, err := utils.ExecuteRequest(c.httpClient, req)
	if err != nil {
		return nil, statusCode, err
	}
	
	return utils.UnmarshalJSONResponse[PurchaseDetailsResponse](body, statusCode, "booking service")
}

