// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    iin,
    username,
    hashed_password,
    name,
    surname,
    email
) VALUES (
             $1, $2, $3, $4, $5, $6
         ) RETURNING iin, username, hashed_password, name, surname, email, created_at
`

type CreateUserParams struct {
	Iin            int64  `json:"iin"`
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Iin,
		arg.Username,
		arg.HashedPassword,
		arg.Name,
		arg.Surname,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.Iin,
		&i.Username,
		&i.HashedPassword,
		&i.Name,
		&i.Surname,
		&i.Email,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE iin = $1
`

func (q *Queries) DeleteUser(ctx context.Context, iin int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, iin)
	return err
}

const getUser = `-- name: GetUser :one
SELECT iin, username, hashed_password, name, surname, email, created_at FROM users
WHERE iin = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, iin int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, iin)
	var i User
	err := row.Scan(
		&i.Iin,
		&i.Username,
		&i.HashedPassword,
		&i.Name,
		&i.Surname,
		&i.Email,
		&i.CreatedAt,
	)
	return i, err
}

const getUsersForUpdate = `-- name: GetUsersForUpdate :one
SELECT iin, username, hashed_password, name, surname, email, created_at FROM users
WHERE iin = $1 LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetUsersForUpdate(ctx context.Context, iin int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUsersForUpdate, iin)
	var i User
	err := row.Scan(
		&i.Iin,
		&i.Username,
		&i.HashedPassword,
		&i.Name,
		&i.Surname,
		&i.Email,
		&i.CreatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT iin, username, hashed_password, name, surname, email, created_at FROM users
ORDER BY iin
    LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.Iin,
			&i.Username,
			&i.HashedPassword,
			&i.Name,
			&i.Surname,
			&i.Email,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2
WHERE iin = $1
    RETURNING iin, username, hashed_password, name, surname, email, created_at
`

type UpdateUserPasswordParams struct {
	Iin            int64  `json:"iin"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPassword, arg.Iin, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.Iin,
		&i.Username,
		&i.HashedPassword,
		&i.Name,
		&i.Surname,
		&i.Email,
		&i.CreatedAt,
	)
	return i, err
}
