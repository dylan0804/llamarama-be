package utils

import (
	"sync"

	db "github.com/dylan0804/Llamarama/cmd/internal/db/sqlc"
	"github.com/dylan0804/Llamarama/cmd/internal/models"
	"github.com/dylan0804/Llamarama/cmd/internal/services"
)

var rooms = make(map[string]*models.Room)
var roomsMutex = &sync.Mutex{}

func GetRoom(roomID string, queries *db.Queries) *models.Room {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	if room, exists := rooms[roomID]; exists {
		return room
	}

	room := &models.Room{
		Clients: make(map[*models.Client]bool),
		Broadcast: make(chan models.Message, 100),
		Mutex: &sync.Mutex{},
	}

	rooms[roomID] = room

	go services.HandleMessages(room, queries)

	return room
}

func AddClient(room *models.Room, client *models.Client) {
	room.Mutex.Lock()
	room.Clients[client] = true
	room.Mutex.Unlock()
}

func RemoveClient(room *models.Room, client *models.Client) {
	room.Mutex.Lock()
	delete(room.Clients, client)
	room.Mutex.Unlock()
}