package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"social/models"
	"social/pkg/db/sqlite"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connection
type Client struct {
	userID string
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
}

// Message represents a WebSocket message for routing
type Message struct {
	Type      string          `json:"type"`
	FromUserID string         `json:"from_user_id,omitempty"`
	ToUserID   string         `json:"to_user_id,omitempty"`
	ToGroupID  string         `json:"to_group_id,omitempty"`
	Payload   json.RawMessage `json:"payload"`
}

// Hub manages WebSocket connections and message routing
type Hub struct {
	clients    map[string]*Client
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// upgrader is used to upgrade HTTP connections to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now, can be restricted later
	},
}

// NewHub creates and returns a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the main event loop for the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			log.Printf("Client registered: %s", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
				client.conn.Close()
				log.Printf("Client unregistered: %s", client.userID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			switch message.Type {
			case "private_message":
				h.SendToUser(message.ToUserID, message)
			case "group_message":
				h.BroadcastToGroup(message.ToGroupID, message)
			case "notification":
				h.SendToUser(message.ToUserID, message)
			}
		}
	}
}

// RegisterClient validates the session token and registers a new client
func (h *Hub) RegisterClient(userID string, conn *websocket.Conn, sessionToken string) error {
	// Validate session token
	session, err := models.GetSessionByID(sqlite.DB, sessionToken)
	if err != nil {
		return errors.New("invalid session")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		models.DeleteSession(sqlite.DB, session.ID)
		return errors.New("session expired")
	}

	// Verify the session belongs to the user
	if session.UserID != userID {
		return errors.New("session does not belong to user")
	}

	// Create and register client
	client := &Client{
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, 256),
		hub:    h,
	}

	h.register <- client

	// Start write pump
	go client.writePump()

	// Start read pump
	go client.readPump()

	return nil
}

// UnregisterClient removes a client from the hub
func (h *Hub) UnregisterClient(userID string) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if ok {
		h.unregister <- client
	}
}

// SendToUser sends a message to a specific user if they are connected
// If not connected, it should save to database (to be implemented in 2.2)
func (h *Hub) SendToUser(userID string, message Message) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if ok {
		// Marshal the message to JSON
		data, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling message for user %s: %v", userID, err)
			return
		}

		// Send to client's send channel (non-blocking)
		select {
		case client.send <- data:
		default:
			// Channel full, client may be disconnected
			log.Printf("Send channel full for user %s, unregistering", userID)
			h.unregister <- client
		}
	} else {
		// User not connected, save to database (placeholder for 2.2)
		h.saveMessageToDatabase(message)
	}
}

// BroadcastToGroup sends a message to all members of a group
func (h *Hub) BroadcastToGroup(groupID string, message Message) {
	// Get group members (placeholder - to be implemented in 2.2)
	members := h.getGroupMembers(groupID)

	for _, userID := range members {
		h.SendToUser(userID, message)
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(51200) // 50KB max message size
	c.conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Client %s normal disconnect: %v", c.userID, err)
			} else {
				log.Printf("Client %s read error: %v", c.userID, err)
			}
			return
		}

		// Parse and validate the incoming message
		var message Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("Invalid message format from client %s: %v", c.userID, err)
			continue
		}

		// Validate message has required fields
		if message.Type == "" {
			log.Printf("Message missing type field from client %s", c.userID)
			continue
		}

		// Set the sender's user ID
		message.FromUserID = c.userID

		// Send to broadcast channel
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Channel closed, hub shut down
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the same frame
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// Send ping to keep connection alive
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// saveMessageToDatabase saves a message to the database when user is offline
// This is a placeholder that will be implemented in sub-issue 2.2
func (h *Hub) saveMessageToDatabase(message Message) {
	// TODO: Implement in sub-issue 2.2
	// This will save the message to the database for offline users
	log.Printf("Saving message to database for offline user: type=%s, to_user_id=%s, to_group_id=%s",
		message.Type, message.ToUserID, message.ToGroupID)
}

// getGroupMembers returns the list of user IDs for a group
// This is a placeholder that will be implemented in sub-issue 2.2
func (h *Hub) getGroupMembers(groupID string) []string {
	// TODO: Implement in sub-issue 2.2
	// This will query the database for group members
	log.Printf("Getting group members for group: %s", groupID)
	return []string{}
}

// ServeWS upgrades an HTTP connection to WebSocket and registers the client
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Get session token from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		conn.Close()
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Register the client
	if err := h.RegisterClient(userID, conn, cookie.Value); err != nil {
		conn.Close()
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// IsUserConnected checks if a user is currently connected
func (h *Hub) IsUserConnected(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// contextKey is used for context values
type contextKey string

// UserIDKey is the context key for user ID
const UserIDKey contextKey = "user_id"

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}