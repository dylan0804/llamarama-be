-- name: CreateUser :one
INSERT INTO users (email, password)
VALUES ($1, $2)
RETURNING id;

-- name: GetUserByEmail :one
SELECT id, password FROM users WHERE email = $1;