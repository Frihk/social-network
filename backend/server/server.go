package server

import (
	"log"
	"net/http"

	"social/pkg/handlers"
	"social/queries/middleware"
	"social/queries/websocket"

	"github.com/gorilla/mux"
)

// WSHub is a global WebSocket hub accessible to other handlers
var WSHub *websocket.Hub

func NewServer() *http.Server {
	router := mux.NewRouter()

	// Apply CORS middleware globally
	router.Use(middleware.CORSMiddleware)

	// Create and start WebSocket hub
	WSHub = websocket.NewHub()
	go WSHub.Run()

	// WebSocket route (NOT protected by auth middleware - session token in URL)
	router.HandleFunc("/api/ws", handlers.HandleWebSocketUpgrade(WSHub)).Methods("GET", "OPTIONS")

	// Public routes
	router.HandleFunc("/api/auth/register", handlers.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/login", handlers.Login).Methods("POST", "OPTIONS")

	// Protected routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/auth/logout", handlers.Logout).Methods("POST", "OPTIONS")
	protected.HandleFunc("/auth/me", handlers.GetMe).Methods("GET", "OPTIONS")
	protected.HandleFunc("/auth/session", handlers.GetSession).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/{id}", handlers.GetUserProfile).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/{id}/privacy", handlers.UpdateProfilePrivacy).Methods("PUT", "OPTIONS")

	// Chat routes
	protected.HandleFunc("/chat/users", handlers.GetDMEligibleUsers).Methods("GET", "OPTIONS")
	protected.HandleFunc("/chat/{userId}", handlers.GetPrivateMessageHistory).Methods("GET", "OPTIONS")
	protected.HandleFunc("/groups/{groupId}/messages", handlers.GetGroupMessageHistory).Methods("GET", "OPTIONS")
	protected.HandleFunc("/user/groups", handlers.GetUserGroups).Methods("GET", "OPTIONS")

	// Notification routes
	protected.HandleFunc("/notifications", handlers.GetNotifications).Methods("GET", "OPTIONS")
	protected.HandleFunc("/notifications/{notificationId}/read", handlers.MarkNotificationAsRead).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/notifications/read-all", handlers.MarkAllNotificationsAsRead).Methods("PUT", "OPTIONS")

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server configured on :8080")
	return server
}
