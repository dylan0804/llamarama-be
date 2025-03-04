// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: rooms.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createRoom = `-- name: CreateRoom :exec
INSERT INTO rooms (name, description)
VALUES ($1, $2)
`

type CreateRoomParams struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) error {
	_, err := q.db.Exec(ctx, createRoom, arg.Name, arg.Description)
	return err
}

const getAllRooms = `-- name: GetAllRooms :many
SELECT id, name, description, created_at, updated_at, is_active FROM rooms
`

func (q *Queries) GetAllRooms(ctx context.Context) ([]Room, error) {
	rows, err := q.db.Query(ctx, getAllRooms)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Room{}
	for rows.Next() {
		var i Room
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.IsActive,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMessagesByRoomID = `-- name: GetMessagesByRoomID :many
SELECT id, user_id, content, created_at FROM messages WHERE room_id = $1
`

type GetMessagesByRoomIDRow struct {
	ID        pgtype.UUID      `json:"id"`
	UserID    pgtype.UUID      `json:"userId"`
	Content   string           `json:"content"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
}

func (q *Queries) GetMessagesByRoomID(ctx context.Context, roomID pgtype.UUID) ([]GetMessagesByRoomIDRow, error) {
	rows, err := q.db.Query(ctx, getMessagesByRoomID, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetMessagesByRoomIDRow{}
	for rows.Next() {
		var i GetMessagesByRoomIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Content,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoomByID = `-- name: GetRoomByID :one
SELECT name, description FROM rooms WHERE id = $1
`

type GetRoomByIDRow struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) GetRoomByID(ctx context.Context, id pgtype.UUID) (GetRoomByIDRow, error) {
	row := q.db.QueryRow(ctx, getRoomByID, id)
	var i GetRoomByIDRow
	err := row.Scan(&i.Name, &i.Description)
	return i, err
}
