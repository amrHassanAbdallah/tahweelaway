-- name: CreateBank :one
INSERT INTO banks (id,
                       name,
                       user_id,
                       branch_number,
                       account_number,
                       account_holder_name,
                       reference,
                       currency,
                       expire_at
                   )
VALUES ($1, $2, $3,$4,$5,$6,$7,$8,$9)
RETURNING *;

-- name: GetUserBanks :one
SELECT *
FROM banks
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
