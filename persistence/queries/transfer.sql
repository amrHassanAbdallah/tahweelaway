-- name: CreateTransfer :one
INSERT INTO transfers (
                       id,
                       from_id,
                       to_id,
                       amount,
                       type)
VALUES ($1, $2, $3,$4,$5)
RETURNING *;

-- name: GetTransfer :one
SELECT *
FROM transfers
WHERE id = $1
LIMIT 1;

-- name: ListTransfers :many
SELECT *
FROM transfers
WHERE from_id = $1
   OR to_id = $2
ORDER BY created_at
LIMIT $3 OFFSET $4;