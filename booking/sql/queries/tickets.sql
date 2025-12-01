-- name: GetTicketWithPrice :one
SELECT 
    t.id,
    t.event_id,
    t.ticket_type_id,
    t.status,
    tt.price_cents
FROM tickets t
JOIN ticket_types tt ON t.ticket_type_id = tt.id
WHERE t.id = $1;

-- name: GetTicketStatus :one
SELECT status FROM tickets WHERE id = $1;

-- name: PurchaseTicket :exec
UPDATE tickets
SET status = 'sold'
WHERE id = $1 AND status = 'available';


