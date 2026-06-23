-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

