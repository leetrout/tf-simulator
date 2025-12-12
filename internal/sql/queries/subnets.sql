-- name: GetSubnet :one
SELECT * FROM subnets
WHERE id = ? LIMIT 1;

-- name: ListSubnets :many
SELECT * FROM subnets;

-- name: CreateSubnet :one
INSERT INTO subnets (name, network_id) VALUES (?, ?)
RETURNING *;
