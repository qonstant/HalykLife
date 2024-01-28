-- name: CreateUser :one
INSERT INTO users (
    iin,
    username,
    hashed_password,
    name,
    surname,
    email,
    user_role
) VALUES (
             $1, $2, $3, $4, $5, $6, $7
         ) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE iin = $1 LIMIT 1;

-- name: GetUsersForUpdate :one
SELECT * FROM users
WHERE iin = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY iin
    LIMIT $1
OFFSET $2;

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2
WHERE iin = $1
    RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE iin = $1;

-- name: GetRoleByUsername :one
SELECT user_role FROM users
WHERE username = $1 LIMIT 1;