-- name: GetServer :one
SELECT * FROM servers
WHERE id = ? LIMIT 1;

-- name: ListServers :many
SELECT * FROM servers;

-- name: CreateServer :one
INSERT INTO servers (name, subnet_id, static_ip_id) VALUES (?, ?, ?)
RETURNING *;
