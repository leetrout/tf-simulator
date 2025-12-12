-- name: GetDatabase :one
SELECT * FROM databases
WHERE id = ? LIMIT 1;

-- name: ListDatabases :many
SELECT * FROM databases;

-- name: CreateDatabase :one
INSERT INTO databases (name, subnet_id) VALUES (?, ?)
RETURNING *;
