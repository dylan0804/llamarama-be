package utils

import (
	"sync"

	"github.com/dylan0804/Llamarama/cmd/internal/models"
	"github.com/dylan0804/Llamarama/cmd/internal/wsc"
)

var rooms = make(map[string]*models.Room)
var roomsMutex = &sync.Mutex{}

func GetRoom(roomID string) *models.Room {
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

	go wsc.HandleMessages(room)

	return room
}