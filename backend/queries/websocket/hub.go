package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"social/models"
	"social/pkg/db/sqlite"
	"social/queries"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients by user ID
	clients map[string]*Client
	// Inbound messages from clients
	broadcast chan []byte
	// Register requests from clients
	register chan *Client
	// Unregister requests from clients
	unregister chan *Client
	// Mutex for thread-safe access
	mutex sync.RWMutex
}

// Client represents a WebSocket client
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
}

// Message represents a WebSocket message
type Message struct {
	Type       string      `json:"type"`
	UserID     string      `json:"user_id,omitempty"`
	GroupID    string      `json:"group_id,omitempty"`
	ID         string      `json:"id,omitempty"`
	SenderID   string      `json:"sender_id,omitempty"`
	ReceiverID string      `json:"receiver_id,omitempty"`
	Content    string      `json:"content,omitempty"`
	CreatedAt  string      `json:"created_at,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.userID] = client
			h.mutex.Unlock()
			log.Printf("Client registered: %s", client.userID)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
				log.Printf("Client unregistered: %s", client.userID)
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			// Handle broadcast messages
			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("Error unmarshaling broadcast message: %v", err)
				continue
			}

			switch msg.Type {
			case "private_message":
				h.SendToUser(msg.UserID, message)
			case "group_message":
				h.BroadcastToGroup(msg.GroupID, message)
			}
		}
	}
}

// RegisterClient validates the session and registers a new client
func (h *Hub) RegisterClient(conn *websocket.Conn, sessionToken string) error {
	session, err := models.GetSessionByID(sqlite.DB, sessionToken)
	if err != nil {
		return err
	}

	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: session.UserID,
	}

	h.register <- client

	// Start read and write pumps
	go client.readPump()
	go client.writePump()

	return nil
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID string, message []byte) {
	h.mutex.RLock()
	client, ok := h.clients[userID]
	h.mutex.RUnlock()

	if ok {
		select {
		case client.send <- message:
		default:
			// Client's send buffer is full, remove them
			h.unregister <- client
		}
	}
	// If user is not connected, message is saved to DB by the caller
}

// BroadcastToGroup broadcasts a message to all members of a group
func (h *Hub) BroadcastToGroup(groupID string, message []byte) {
	// Query all accepted group members
	rows, err := sqlite.DB.Query(`
		SELECT user_id FROM group_members 
		WHERE group_id = ? AND status = 'accepted'
	`, groupID)
	if err != nil {
		log.Printf("Error querying group members: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			log.Printf("Error scanning user_id: %v", err)
			continue
		}
		h.SendToUser(userID, message)
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		// Handle incoming messages
		var msg Message
		if err := json.Unmarshal(message, &msg); err == nil {
			msg.SenderID = c.userID
			msg.CreatedAt = time.Now().UTC().Format(time.RFC3339)

			if msg.Type == "private_message" && msg.UserID != "" {
				id, err := queries.SavePrivateMessage(c.userID, msg.UserID, msg.Content)
				if err != nil {
					log.Printf("Error saving private message: %v", err)
				} else {
					msg.ID = id
					msg.ReceiverID = msg.UserID
				}
			} else if msg.Type == "group_message" && msg.GroupID != "" {
				id, err := queries.SaveGroupMessage(msg.GroupID, c.userID, msg.Content)
				if err != nil {
					log.Printf("Error saving group message: %v", err)
				} else {
					msg.ID = id
				}
			}

			if updatedMsg, err := json.Marshal(msg); err == nil {
				c.hub.broadcast <- updatedMsg
			} else {
				c.hub.broadcast <- message
			}
		} else {
			c.hub.broadcast <- message
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}
		}
	}
}
