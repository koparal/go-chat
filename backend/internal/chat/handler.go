package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	chat *Chat
}

func NewHandler(chat *Chat) *Handler {
	return &Handler{
		chat: chat,
	}
}

// CreateRoom godoc
// @Summary Create a new chat room
// @Description Create a new chat room with the provided data
// @Tags rooms
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header"
// @Param input body Room true "Room data to create"
// @Success 200 {object} Room "OK"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rooms/create [post]
func (h *Handler) CreateRoom(c *gin.Context) {
	var room Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomKey := fmt.Sprintf("room:%s", room.ID)
	roomJSON, _ := json.Marshal(room)

	ctx := context.Background()
	if err := h.chat.redis.Set(ctx, roomKey, roomJSON, 0).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	c.JSON(http.StatusOK, room)
}

// GetRooms godoc
// @Summary Get all chat rooms
// @Description Get a list of all available chat rooms
// @Tags rooms
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} []Room "OK"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rooms/list [get]
func (h *Handler) GetRooms(c *gin.Context) {
	ctx := context.Background()
	keys := h.chat.redis.Keys(ctx, "room:*").Val()
	rooms := make([]Room, 0)

	for _, key := range keys {
		roomJSON, err := h.chat.redis.Get(ctx, key).Bytes()
		if err != nil {
			log.Printf("Error retrieving room data for key %s: %v", key, err)
			continue
		}

		var room Room
		if err := json.Unmarshal(roomJSON, &room); err != nil {
			log.Printf("Error unmarshalling room data for key %s: %v", key, err)
			continue
		}

		rooms = append(rooms, room)
	}

	c.JSON(http.StatusOK, rooms)
}

// JoinRoom godoc
// @Summary Join a chat room
// @Description Join a chat room with the provided room ID, user ID, and username
// @Tags rooms
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header"
// @Param roomId path string true "Room ID"
// @Param userId query string true "User ID"
// @Param username query string true "Username"
// @Success 200 {object} Message "OK"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rooms/join/{roomId} [get]
func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	h.chat.Register <- cl
	defer func() { h.chat.Unregister <- cl }()

	m := &Message{
		Content:  "New user has joined.",
		RoomID:   roomID,
		Username: username,
	}
	h.chat.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.chat)
}

// GetClients godoc
// @Summary Get clients in a chat room
// @Description Get the list of clients currently in a chat room
// @Tags rooms
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header"
// @Param roomId path string true "Room ID"
// @Success 200 {object} []Client "OK"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rooms/clients/{roomId} [get]
func (h *Handler) GetClients(c *gin.Context) {
	roomId := c.Param("roomId")
	roomKey := fmt.Sprintf("room:%s", roomId)
	fmt.Println(roomKey)
	ctx := context.Background()
	roomJSON, err := h.chat.redis.Get(ctx, roomKey).Bytes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get clients"})
		return
	}

	var room Room
	if err := json.Unmarshal(roomJSON, &room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal room data"})
		return
	}

	var clients []Client
	for _, cl := range room.Clients {
		clients = append(clients, *cl)
	}

	c.JSON(http.StatusOK, clients)
}
