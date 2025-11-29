-- name: CreateVenue :one
INSERT INTO venues (name, location, seat_map)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetVenue :one
SELECT * FROM venues
WHERE id = $1;

-- name: GetVenues :many
SELECT * FROM venues
ORDER BY name ASC
LIMIT $1
OFFSET $2;

-- name: UpdateVenue :one
UPDATE venues
SET name = $2,
    location = $3,
    seat_map = $4
WHERE id = $1
RETURNING *;

-- name: DeleteVenue :exec
DELETE FROM venues
WHERE id = $1;
