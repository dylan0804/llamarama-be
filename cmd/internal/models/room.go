package models

import (
	"sync"
	"time"

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
	UserID string `json:"user_id"`
}

type Message struct {
	Sender *Client
	Payload MessagePayload
}

type Room struct {
	ID string
	Clients map[*Client]bool
	Broadcast chan Message
	Mutex *sync.Mutex
}

type RoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoomResponse struct {
	RoomDetails RoomDetails     `json:"roomDetails"`
	Messages    []MessageDetail `json:"messages"` 
}

type RoomDetails struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MessageDetail struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}