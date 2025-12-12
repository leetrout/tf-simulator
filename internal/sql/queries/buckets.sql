-- name: GetBucket :one
SELECT * FROM buckets
WHERE id = ? LIMIT 1;

-- name: ListBuckets :many
SELECT * FROM buckets;

-- name: CreateBucket :one
INSERT INTO buckets (name) VALUES (?)
RETURNING *;
