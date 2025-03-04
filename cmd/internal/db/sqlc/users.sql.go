// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, password)
VALUES ($1, $2)
RETURNING id
`

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Email, arg.Password)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, password FROM users WHERE email = $1
`

type GetUserByEmailRow struct {
	ID       pgtype.UUID `json:"id"`
	Password string      `json:"password"`
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(&i.ID, &i.Password)
	return i, err
}
