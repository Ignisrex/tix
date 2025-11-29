-- name: CreateEvent :one
INSERT INTO events (title, description, start_date, venue_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetEvent :one
SELECT * FROM events
WHERE id = $1;

-- name: GetEvents :many
SELECT * FROM events
ORDER BY start_date DESC
LIMIT $1
OFFSET $2;

-- name: UpdateEvent :one
UPDATE events SET title = $2, description = $3, start_date = $4, venue_id = $5
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1;