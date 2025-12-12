-- name: GetLoadBalancer :one
SELECT * FROM load_balancers
WHERE id = ? LIMIT 1;

-- name: ListLoadBalancers :many
SELECT * FROM load_balancers;

-- name: CreateLoadBalancer :one
INSERT INTO load_balancers (name, subnet_id, static_ip_id) VALUES (?, ?, ?)
RETURNING *;
