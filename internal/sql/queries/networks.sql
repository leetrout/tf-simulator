-- name: GetNetwork :one
SELECT * FROM networks
WHERE id = ? LIMIT 1;

-- name: ListNetworks :many
SELECT * FROM networks;

-- name: CreateNetwork :one
INSERT INTO networks (name) VALUES (?)
RETURNING *;
