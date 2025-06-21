-- name: CreateOneUser :one
INSERT INTO users
(id, email, role, password_hash)
VALUES
($1, $2, $3, $4)
RETURNING id;

-- name: GetOneUser :one
SELECT
    id,
    email,
    role,
    created_at,
    updated_at
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetOneCredentialUserByEmail :one
SELECT
    id,
    password_hash
FROM users
WHERE
    email = $1 AND deleted_at IS NULL;

-- name: GetManyUser :many
SELECT
    id,
    email,
    role
FROM users
WHERE
    deleted_at IS NULL
LIMIT $1 OFFSET $2;

-- name: UpdateOnePasswordUser :one
UPDATE users
SET password_hash = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING id;

-- name: UpdateOneEmailUser :one
UPDATE users
SET email = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING id;

-- name: UpdateOneRoleUser :one
UPDATE users
SET role = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING id;

-- name: DeleteSoftOneUser :one
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING id;

-- name: DeleteHardOneUser :one
DELETE FROM users
WHERE id = $1
RETURNING id;

-- name: CountEmailUser :one
SELECT COUNT(*) FROM users
WHERE email = $1;

-- name: CountIDUser :one
SELECT COUNT(*) FROM users
WHERE id = $1;
