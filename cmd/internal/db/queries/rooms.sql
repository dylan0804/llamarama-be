-- name: GetAllRooms :many
SELECT * FROM rooms;

-- name: GetRoomByID :one
SELECT name, description FROM rooms WHERE id = $1;

-- name: GetMessagesByRoomID :many
SELECT id, user_id, content, created_at FROM messages WHERE room_id = $1;

-- name: CreateRoom :exec
INSERT INTO rooms (name, description)
VALUES ($1, $2);