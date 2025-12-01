package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"

	"github.com/ignisrex/tix/core/types"
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

	client := &Client{es: es}

	if err := client.EnsureIndex(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure index exists: %w", err)
	}

	return client, nil
}


func (c *Client) EnsureIndex(ctx context.Context) error {
	indexName := "events"

	res, err := c.es.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	mapping := `{
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"title": { 
					"type": "text",
					"analyzer": "standard",
					"fields": {
						"keyword": { "type": "keyword" }
					}
				},
				"description": { 
					"type": "text",
					"analyzer": "standard"
				},
				"start_date": { "type": "date" },
				"venue_id": { "type": "keyword" },
				"venue_name": { 
					"type": "text",
					"fields": {
						"keyword": { "type": "keyword" }
					}
				},
				"venue_location": { 
					"type": "text",
					"fields": {
						"keyword": { "type": "keyword" }
					}
				},
				"created_at": { "type": "date" }
			}
		},
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0,
			"refresh_interval": "1s"
		}
	}`

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	res, err = req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}

	return nil
}


func (c *Client) IndexEvent(ctx context.Context, event types.Event, venue types.Venue) error {
	indexName := "events"

	doc := map[string]interface{}{
		"id":             event.ID.String(),
		"title":          event.Title,
		"description":    event.Description,
		"start_date":     event.StartDate.Format("2006-01-02T15:04:05Z07:00"),
		"venue_id":       event.VenueID.String(),
		"venue_name":      venue.Name,
		"venue_location":  venue.Location,
		"created_at":      event.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), 
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: event.ID.String(),
		Body:       strings.NewReader(string(docJSON)),
		Refresh:    "true", 
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}


func (c *Client) UpdateEvent(ctx context.Context, event types.Event, venue types.Venue) error {
	return c.IndexEvent(ctx, event, venue)
}

func (c *Client) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	indexName := "events"

	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: eventID.String(),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting document: %s", res.String())
	}

	return nil
}

