package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"

	db "github.com/dylan0804/Llamarama/cmd/internal/db/sqlc"
	"github.com/dylan0804/Llamarama/cmd/internal/models"
	"github.com/dylan0804/Llamarama/cmd/internal/response"
	"github.com/dylan0804/Llamarama/cmd/internal/services"
	"github.com/dylan0804/Llamarama/cmd/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	queries *db.Queries
	sessionStore *utils.SessionStore
}

func NewHandler(queries *db.Queries, sessionStore *utils.SessionStore) *Handler {
	return &Handler{
		queries: queries,
		sessionStore: sessionStore,
	}
}

var upgrade = websocket.Upgrader{
	CheckOrigin: func (r *http.Request) bool {
		return true
	},
}

func (h *Handler) WsHandler(c *gin.Context) {	
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()

	userId := c.MustGet("user_id").(string)
	
	client := &models.Client{
		ID: userId,
		Conn: conn,
	}
	
	roomID := c.Param("id")

	room := utils.GetRoom(roomID, h.queries)

	room.ID = roomID

	utils.AddClient(room, client)

	services.ReadMessages(room, client)
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

	var room db.GetRoomByIDRow
	var messages []db.GetMessagesByRoomIdRow
	var roomErr, msgError error

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		room, roomErr = h.queries.GetRoomByID(c.Request.Context(), pgtype.UUID{
			Bytes: id,
			Valid: true,
		})
	}()

	go func() {
		defer wg.Done()
		messages, msgError = h.queries.GetMessagesByRoomId(c.Request.Context(), pgtype.UUID{
			Bytes: id,
			Valid: true,
		})
	}()

	wg.Wait()

	if roomErr != nil {
		response.Error(c.Writer, http.StatusInternalServerError, "Error getting room by ID", roomErr.Error())
		return
	}
	if msgError != nil {
		response.Error(c.Writer, http.StatusInternalServerError, "Error getting messages by room ID", msgError.Error())
		return
	}

	response.Success(c.Writer, http.StatusOK, "Room fetched successfully", gin.H{
		"roomDetails": room,
		"messages": messages,
	})
}

func (h *Handler) Register(c *gin.Context) {
	var userReq models.UserRequest

	if err := c.ShouldBindJSON(&userReq); err != nil {
		response.Error(c.Writer, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c.Writer, http.StatusInternalServerError, "Failed to hash password", err.Error())
		return
	}

	id, err := h.queries.CreateUser(c.Request.Context(), db.CreateUserParams{
		Email: userReq.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            if strings.Contains(pgErr.ConstraintName, "email") {
				log.Println("error", pgErr.ConstraintName)
                response.Error(c.Writer, http.StatusConflict, "Email already exists", "A user with this email already exists")
                return
            }
        }

		response.Error(c.Writer, http.StatusInternalServerError, "Failed to create user", err.Error())
        return
    }

	token, err := h.sessionStore.CreateToken(c.Request.Context(), id.String())
	if err != nil {
		response.Error(c.Writer, http.StatusInternalServerError, "Failed to create token", err.Error())
		return
	}

	response.Success(c.Writer, http.StatusCreated, "User created successfully", gin.H{
		"token": token,
		"user_id": id.String(),
	})
}

func (h *Handler) Login(c *gin.Context) {
	var userReq models.UserRequest

	if err := c.ShouldBindJSON(&userReq); err != nil {
		response.Error(c.Writer, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.queries.GetUserByEmail(c.Request.Context(), userReq.Email)
	if err != nil {
		response.Error(c.Writer, http.StatusInternalServerError, "Invalid email or password", err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userReq.Password)); err != nil {
		response.Error(c.Writer, http.StatusUnauthorized, "Invalid password", err.Error())
		return
	}

	token, err := h.sessionStore.CreateToken(c.Request.Context(), user.ID.String())
	if err != nil {
		response.Error(c.Writer, http.StatusInternalServerError, "Failed to create token", err.Error())
		return
	}

	response.Success(c.Writer, http.StatusOK, "Login successful", gin.H{
		"token": token,
		"user_id": user.ID.String(),
	})
}

func (h *Handler) Logout(c *gin.Context) {
	h.sessionStore.Delete(c.Request.Context(), c.MustGet("user_token").(string))

	response.Success(c.Writer, http.StatusOK, "Logout successful", nil)
}