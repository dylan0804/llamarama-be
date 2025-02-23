package wsc

import (
	"fmt"
	"net/http"

	"github.com/dylan0804/Llamarama/cmd/internal/models"
	"github.com/gorilla/websocket"
)

var Upgrade = websocket.Upgrader{
	CheckOrigin: func (r *http.Request) bool {
		return true
	},
}

func HandleMessages(room *models.Room) {
	for {
		message := <- room.Broadcast
		fmt.Printf("Handling broadcast message\n") // Debug log

		room.Mutex.Lock()
		for client := range room.Clients {
			if client == message.Sender {
				fmt.Println("Skipping message to self")
				continue
			}

			fmt.Println("Sending message to client", message.Payload.Content)
			err := client.Conn.WriteMessage(websocket.TextMessage, message.Payload.Content)
			if err != nil {
				client.Conn.Close()
				delete(room.Clients, client)
			}
		}
		room.Mutex.Unlock()
	}
}