-- name: BatchCreateTickets :many
INSERT INTO tickets (event_id, ticket_type_id, status)
SELECT 
    $1::uuid,
    unnest($2::uuid[]),
    'available'::ticket_status
RETURNING id, event_id, ticket_type_id, status, created_at, updated_at;

-- name: GetTicketsForEvent :many
SELECT * FROM tickets
WHERE event_id = $1;

-- name: GetTicket :one
SELECT * FROM tickets
WHERE event_id = $1 AND id = $2;


-- name: GetTicketTypes :many
SELECT * FROM ticket_types;