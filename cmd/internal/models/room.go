package models

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID string
	Email string
	Conn *websocket.Conn
}

type MessagePayload struct {
	Type string `json:"type"`
	Content string `json:"content"`
	ClientID string `json:"client_id"`
}

type Message struct {
	Sender *Client
	Payload json.RawMessage
}

type Room struct {
	ID string
	Clients map[*Client]bool
	Broadcast chan Message
	Mutex *sync.Mutex
}

type RoomRequest struct {
	Name string `json:"name"`
	Description string `json:"description"`
}