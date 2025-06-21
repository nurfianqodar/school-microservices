-- name: CreateUser :one
INSERT INTO users
(id, email, role, password_hash)
VALUES
($1, $2, $3, $4)
RETURNING id;
