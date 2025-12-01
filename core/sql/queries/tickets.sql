-- name: BatchCreateTickets :many
INSERT INTO tickets (event_id, ticket_type_id, status)
SELECT 
    $1::uuid,
    unnest($2::uuid[]),
    'available'::ticket_status
RETURNING id, event_id, ticket_type_id, status, created_at, updated_at;

-- name: GetTicketsForEvent :many
SELECT 
    id,
    event_id,
    ticket_type_id,
    status,
    created_at,
    updated_at,
    ticket_type_name,
    ticket_type_display_name,
    ticket_type_price_cents
FROM enriched_tickets
WHERE event_id = $1
ORDER BY ticket_type_id, id;

-- name: GetTicket :one
SELECT 
    id,
    event_id,
    ticket_type_id,
    status,
    created_at,
    updated_at,
    ticket_type_name,
    ticket_type_display_name,
    ticket_type_price_cents
FROM enriched_tickets
WHERE event_id = $1 AND id = $2;


-- name: GetTicketTypes :many
SELECT * FROM ticket_types;