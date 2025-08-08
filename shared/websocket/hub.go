package websocket

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"ua/shared/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Hub struct {
	clients    map[uuid.UUID]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mutex      sync.RWMutex
	gameRooms  map[uuid.UUID]*GameRoom
}

type Client struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
	GameID *uuid.UUID
}

type GameRoom struct {
	ID      uuid.UUID
	Clients map[uuid.UUID]*Client
	mutex   sync.RWMutex
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	From    *uuid.UUID  `json:"from,omitempty"`
	To      *uuid.UUID  `json:"to,omitempty"`
	GameID  *uuid.UUID  `json:"game_id,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		gameRooms:  make(map[uuid.UUID]*GameRoom),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.ID] = client
			h.mutex.Unlock()
			logger.Info("Client connected", zap.String("client_id", client.ID.String()))

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
				if client.GameID != nil {
					h.removeClientFromGameRoom(*client.GameID, client.ID)
				}
			}
			h.mutex.Unlock()
			logger.Info("Client disconnected", zap.String("client_id", client.ID.String()))

		case message := <-h.broadcast:
			h.mutex.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					delete(h.clients, client.ID)
					close(client.Send)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) HandleWebSocket(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}

	client := &Client{
		ID:     uuid.New(),
		UserID: userID.(uuid.UUID),
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    h,
	}

	h.register <- client

	go client.writePump()
	go client.readPump()
}

func (h *Hub) JoinGameRoom(gameID uuid.UUID, clientID uuid.UUID) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	client, exists := h.clients[clientID]
	if !exists {
		return
	}

	if _, exists := h.gameRooms[gameID]; !exists {
		h.gameRooms[gameID] = &GameRoom{
			ID:      gameID,
			Clients: make(map[uuid.UUID]*Client),
		}
	}

	room := h.gameRooms[gameID]
	room.mutex.Lock()
	room.Clients[clientID] = client
	room.mutex.Unlock()

	client.GameID = &gameID
}

func (h *Hub) removeClientFromGameRoom(gameID uuid.UUID, clientID uuid.UUID) {
	room, exists := h.gameRooms[gameID]
	if !exists {
		return
	}

	room.mutex.Lock()
	delete(room.Clients, clientID)
	isEmpty := len(room.Clients) == 0
	room.mutex.Unlock()

	if isEmpty {
		delete(h.gameRooms, gameID)
	}
}

func (h *Hub) BroadcastToGame(gameID uuid.UUID, message []byte) {
	h.mutex.RLock()
	room, exists := h.gameRooms[gameID]
	h.mutex.RUnlock()

	if !exists {
		return
	}

	room.mutex.RLock()
	defer room.mutex.RUnlock()

	for _, client := range room.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(room.Clients, client.ID)
		}
	}
}

func (h *Hub) SendToUser(userID uuid.UUID, message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for _, client := range h.clients {
		if client.UserID == userID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client.ID)
			}
			return
		}
	}
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("WebSocket error", zap.Error(err))
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			logger.Error("Invalid message format", zap.Error(err))
			continue
		}

		c.handleMessage(&msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(msg *Message) {
}
