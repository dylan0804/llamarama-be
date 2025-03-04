package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	db "github.com/dylan0804/Llamarama/cmd/internal/db/sqlc"
	"github.com/dylan0804/Llamarama/cmd/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgtype"
)

func ReadMessages(room *models.Room, client *models.Client) {
	for {
		_, message, err := client.Conn.ReadMessage()

		if err != nil {
			break
		}

		var msg models.MessagePayload
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		msg.UserID = client.ID

		switch msg.Type {
		case "ping":
			client.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "pong"}`))
		case "message":
			room.Broadcast <- models.Message{
				Sender: client,
				Payload: msg,
			}
		default:
			fmt.Println("Received unknown message type:", msg.Type)
		}
	}
}

func HandleMessages(room *models.Room, queries *db.Queries) {
	for {
		message := <- room.Broadcast

		payload, err := json.Marshal(message.Payload)
		if err != nil {
			log.Println("Error marshalling message:", err)
			continue
		}

		room.Mutex.Lock()
		for client := range room.Clients {
			fmt.Println("Sending message to client", message)
	
			err = client.Conn.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				client.Conn.Close()
				delete(room.Clients, client)
			}
		}

		roomID, err := uuid.Parse(room.ID)
		if err != nil {
			log.Println("Error parsing room ID:", err)
			continue
		}

		userID, err := uuid.Parse(message.Sender.ID)
		if err != nil {
			log.Println("Error parsing user ID:", err)
			continue
		}

		err = queries.CreateMessage(context.Background(), db.CreateMessageParams{
			UserID: pgtype.UUID{
				Bytes: userID,
				Valid: true,
			},
			RoomID: pgtype.UUID{
				Bytes: roomID,
				Valid: true,
			},
			Content: message.Payload.Content,
		})
		if err != nil {
			log.Println("Error inserting message:", err)
			log.Printf("Message content being inserted: %s", string(payload))
			continue
		}

		room.Mutex.Unlock()
	}
}