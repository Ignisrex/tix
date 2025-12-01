package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
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

type ReserveResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	TicketID uuid.UUID `json:"ticket_id"`
}

type PurchaseResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	TicketID uuid.UUID `json:"ticket_id"`
	Total    int32     `json:"total"`
}

func (c *Client) ReserveTicket(ctx context.Context, ticketID uuid.UUID) (*ReserveResponse, int, error) {
	url := fmt.Sprintf("%s/api/v1/booking/reserve/%s", c.baseURL, ticketID.String())
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to call booking service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
	}
	
	var reserveResp ReserveResponse
	if err := json.Unmarshal(body, &reserveResp); err != nil {
		
		if resp.StatusCode != http.StatusOK {
			return nil, resp.StatusCode, fmt.Errorf("booking service returned status %d: %s", resp.StatusCode, string(body))
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to decode response: %w", err)
	}

	return &reserveResp, resp.StatusCode, nil
}

func (c *Client) PurchaseTicket(ctx context.Context, ticketID uuid.UUID) (*PurchaseResponse, int, error) {
	url := fmt.Sprintf("%s/api/v1/booking/purchase/%s", c.baseURL, ticketID.String())
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to call booking service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
	}

	
	var purchaseResp PurchaseResponse
	if err := json.Unmarshal(body, &purchaseResp); err != nil {
		
		if resp.StatusCode != http.StatusOK {
			return nil, resp.StatusCode, fmt.Errorf("booking service returned status %d: %s", resp.StatusCode, string(body))
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to decode response: %w", err)
	}

	return &purchaseResp, resp.StatusCode, nil
}

