package models

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID string
	Conn *websocket.Conn
}

type MessagePayload struct {
	Type string `json:"type"`
	Content json.RawMessage `json:"content"`
}

type Message struct {
	Sender *Client
	Type string
	Payload MessagePayload
}

type Room struct {
	Clients map[*Client]bool
	Broadcast chan Message
	Mutex *sync.Mutex
}

type RoomRequest struct {
	Name string `json:"name"`
	Description string `json:"description"`
}