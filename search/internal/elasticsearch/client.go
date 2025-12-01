package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Client struct {
	es *elasticsearch.Client
}

func NewClient(addresses []string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	return &Client{es: es}, nil
}

type SearchResult struct {
	ID             string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	VenueID         string    `json:"venue_id"`
	VenueName       string    `json:"venue_name"`
	VenueLocation   string    `json:"venue_location"`
	CreatedAt       time.Time `json:"created_at"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}

// Only returns future events (start_date >= now)
func (c *Client) SearchEvents(ctx context.Context, query string, limit, offset int) (*SearchResponse, error) {
	indexName := "events"
	
	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	
	searchQuery := map[string]interface{}{
		"size": limit,
		"from": offset,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  query,
							"fields": []string{"title^2", "description", "venue_name", "venue_location"},
							"type":   "best_fields",
							"fuzziness": "AUTO",
						},
					},
				},
				"filter": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"start_date": map[string]interface{}{
								"gte": now,
							},
						},
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{
				"start_date": map[string]interface{}{
					"order": "asc",
				},
			},
		},
	}

	queryJSON, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(queryJSON),
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	// Parse the response
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return &SearchResponse{Results: []SearchResult{}, Total: 0}, nil
	}

	total, _ := hits["total"].(map[string]interface{})
	totalValue := 0
	if total != nil {
		if val, ok := total["value"].(float64); ok {
			totalValue = int(val)
		}
	}

	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		return &SearchResponse{Results: []SearchResult{}, Total: totalValue}, nil
	}

	results := make([]SearchResult, 0, len(hitsArray))
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		result := SearchResult{
			ID:           getString(source, "id"),
			Title:        getString(source, "title"),
			Description:  getString(source, "description"),
			VenueID:      getString(source, "venue_id"),
			VenueName:    getString(source, "venue_name"),
			VenueLocation: getString(source, "venue_location"),
		}

		// Parse dates
		if startDateStr := getString(source, "start_date"); startDateStr != "" {
			if t, err := time.Parse("2006-01-02T15:04:05Z07:00", startDateStr); err == nil {
				result.StartDate = t
			}
		}
		if createdAtStr := getString(source, "created_at"); createdAtStr != "" {
			if t, err := time.Parse("2006-01-02T15:04:05Z07:00", createdAtStr); err == nil {
				result.CreatedAt = t
			}
		}

		results = append(results, result)
	}

	return &SearchResponse{
		Results: results,
		Total:   totalValue,
	}, nil
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

