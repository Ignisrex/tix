-- name: GetTicketsWithPrice :many
SELECT 
    t.id,
    t.event_id,
    t.ticket_type_id,
    t.status,
    tt.price_cents
FROM tickets t
JOIN ticket_types tt ON t.ticket_type_id = tt.id
WHERE t.id = ANY($1::uuid[]);

-- This query creates a purchase record and updates all tickets atomically
-- name: PurchaseTickets :one
WITH purchase_insert AS (
    INSERT INTO purchases (total_cents)
    VALUES ($1)
    RETURNING id
),
updated_tickets AS (
    UPDATE tickets
    SET status = 'sold', purchase_id = (SELECT id FROM purchase_insert)
    WHERE id = ANY($2::uuid[]) AND status = 'available'
    RETURNING purchase_id
)
SELECT id FROM purchase_insert;

-- need to test this query performance; can converted to view?
-- name: GetPurchaseDetails :one
SELECT 
    p.id as purchase_id,
    p.total_cents,
    p.created_at as purchase_created_at,
    json_agg(
        json_build_object(
            'id', t.id,
            'event_id', t.event_id,
            'ticket_type_id', t.ticket_type_id,
            'status', t.status,
            'ticket_type_name', tt.name,
            'ticket_type_display_name', tt.display_name,
            'ticket_type_price_cents', tt.price_cents
        )
    ) as tickets
FROM purchases p
JOIN tickets t ON t.purchase_id = p.id
JOIN ticket_types tt ON t.ticket_type_id = tt.id
WHERE p.id = $1
GROUP BY p.id, p.total_cents, p.created_at;


