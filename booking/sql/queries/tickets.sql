-- name: ReserveTicket :exec
UPDATE tickets
SET status = 'booked'
WHERE id = $1 AND status = 'available';

-- name: PurchaseTicket :exec
UPDATE tickets
SET status = 'sold'
WHERE id = $1 AND status = 'booked';


