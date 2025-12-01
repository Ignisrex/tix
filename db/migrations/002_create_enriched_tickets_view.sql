-- +goose Up
-- Create view for enriched tickets with ticket type information
CREATE OR REPLACE VIEW enriched_tickets AS
SELECT 
    t.id,
    t.event_id,
    t.ticket_type_id,
    t.status,
    t.created_at,
    t.updated_at,
    tt.name AS ticket_type_name,
    tt.display_name AS ticket_type_display_name,
    tt.price_cents AS ticket_type_price_cents
FROM tickets t
JOIN ticket_types tt ON t.ticket_type_id = tt.id;


-- +goose Down
DROP VIEW IF EXISTS enriched_tickets;

