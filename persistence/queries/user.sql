-- name: CreateUser :one
INSERT INTO users (id,
                   name,
                   username,
                   email,
                   password,
                   currency)
VALUES ($1, $2, $3, $4, $5,$6)
RETURNING *;

-- name: GetUserByUserName :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: AddUserBalance :one
UPDATE users
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeductUserBalance :one
UPDATE users
SET balance = balance - sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;