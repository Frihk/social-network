package handlers

import (
	"net/http"

	ws "social/queries/websocket"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000" || origin == "http://localhost:3001"
	},
}

// HandleWebSocketUpgrade returns a handler function for upgrading HTTP to WebSocket
func HandleWebSocketUpgrade(hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract session token from query parameter or Authorization header
		sessionToken := r.URL.Query().Get("session_token")
		if sessionToken == "" {
			// Try Authorization header
			sessionToken = r.Header.Get("Authorization")
			if len(sessionToken) > 7 && sessionToken[:7] == "Bearer " {
				sessionToken = sessionToken[7:]
			}
		}

		if sessionToken == "" {
			http.Error(w, `{"error":"Missing session token"}`, http.StatusUnauthorized)
			return
		}

		// Upgrade the HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, `{"error":"Failed to upgrade to WebSocket"}`, http.StatusBadRequest)
			return
		}

		// Register the client with the hub
		if err := hub.RegisterClient(conn, sessionToken); err != nil {
			conn.Close()
			http.Error(w, `{"error":"Invalid session token"}`, http.StatusUnauthorized)
			return
		}
	}
}