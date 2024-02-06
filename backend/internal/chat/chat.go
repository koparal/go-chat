package chat

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
)

func NewChat(redis *redis.Client) *Chat {
	return &Chat{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
		redis:      redis,
	}
}

func (c *Chat) Start() {
	for {
		select {
		case cl := <-c.Register:
			c.mu.Lock()
			room, ok := c.Rooms[cl.RoomID]
			if !ok {
				room = &Room{
					ID:      cl.RoomID,
					Name:    "Room " + cl.RoomID,
					Clients: make(map[string]*Client),
				}
				c.Rooms[cl.RoomID] = room
			}
			room.Clients[cl.ID] = cl
			c.mu.Unlock()

		case cl := <-c.Unregister:
			c.mu.Lock()
			if room, ok := c.Rooms[cl.RoomID]; ok {
				delete(room.Clients, cl.ID)
				if len(room.Clients) == 0 {
					delete(c.Rooms, cl.RoomID)
				}
			}
			c.mu.Unlock()

		case m := <-c.Broadcast:
			if room, ok := c.Rooms[m.RoomID]; ok {
				for _, cl := range room.Clients {
					cl.Message <- m
				}
			}
		}
	}
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) readMessage(chat *Chat) {
	defer func() {
		chat.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, 1001, 1006) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}

		chat.Broadcast <- msg
	}
}
