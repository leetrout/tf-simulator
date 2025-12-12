-- name: GetStaticIP :one
SELECT * FROM static_ips
WHERE id = ? LIMIT 1;

-- name: ListStaticIPs :many
SELECT * FROM static_ips;

-- name: CreateStaticIP :one
INSERT INTO static_ips (name, ip_address) VALUES (?, ?)
RETURNING *;
