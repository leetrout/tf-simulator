-- name: GetDomainName :one
SELECT * FROM domain_names
WHERE id = ? LIMIT 1;

-- name: ListDomainNames :many
SELECT * FROM domain_names;

-- name: CreateDomainName :one
INSERT INTO domain_names (name, static_ip_id) VALUES (?, ?)
RETURNING *;
