package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "github.com/dylan0804/Llamarama/cmd/internal/db/sqlc"
	"github.com/dylan0804/Llamarama/cmd/internal/models"
	"github.com/dylan0804/Llamarama/cmd/internal/response"
	"github.com/dylan0804/Llamarama/cmd/internal/utils"
	"github.com/dylan0804/Llamarama/cmd/internal/wsc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	queries *db.Queries
}

func NewHandler(queries *db.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}

func (h *Handler) WsHandler(c *gin.Context) {	
	conn, err := wsc.Upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()

	client := &models.Client{
		ID: uuid.New().String(),
		Conn: conn,
	}
	
	roomID := c.Param("id")

	room := utils.GetRoom(roomID)

	room.Mutex.Lock()
	room.Clients[client] = true
	room.Mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			room.Mutex.Lock()
			delete(room.Clients, client)
			room.Mutex.Unlock()
			break
		}

		var msg models.MessagePayload
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		switch msg.Type {
		case "ping":
			client.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "pong"}`))
		case "message":
			room.Broadcast <- models.Message{
				Sender: client,
				Type: "message",
				Payload: msg,
			}
		default:
			fmt.Println("Received unknown message type:", msg.Type)
		}
	}
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var roomReq models.RoomRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&roomReq); err != nil {
		log.Println("Error decoding room request:", err)
		response.Error(c.Writer, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err := h.queries.CreateRoom(c.Request.Context(), db.CreateRoomParams{
		Name: roomReq.Name,
		Description: pgtype.Text{
			Valid: true,
			String: roomReq.Description,
		},
	})
	if err != nil {
		log.Println("Error creating room:", err)
		response.Error(c.Writer, http.StatusInternalServerError, "Failed to create room", err.Error())
		return
	}

	response.Success(c.Writer, http.StatusCreated, "Room created successfully", nil)
}

func (h *Handler) ListRooms(c *gin.Context) {
	var rooms []db.Room

	rooms, err := h.queries.GetAllRooms(c.Request.Context())
	if err != nil {
		log.Println("Error getting all rooms:", err)
		response.Error(c.Writer, http.StatusInternalServerError, "Failed to get all rooms", err.Error())
		return
	}

	response.Success(c.Writer, http.StatusOK, "Rooms fetched successfully", rooms)
}

func (h *Handler) GetRoom(c *gin.Context) {
	roomID := c.Param("id")

	id, err := uuid.Parse(roomID)
	if err != nil {
		response.Error(c.Writer, http.StatusBadRequest, "Error parsing room ID", err.Error())
		return
	}

	room, err := h.queries.GetRoomById(c.Request.Context(), pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		response.Error(c.Writer, http.StatusInternalServerError, "Error getting room by ID", err.Error())
		return
	}

	response.Success(c.Writer, http.StatusOK, "Room fetched successfully", room)
}