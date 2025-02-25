-- name: CreateMessage :exec
INSERT INTO messages (user_id, room_id, content)
VALUES ($1, $2, $3);

-- name: GetMessagesByRoomId :many
SELECT user_id, room_id, content, created_at FROM messages WHERE room_id = $1;