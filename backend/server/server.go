package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"social/feature"
	"social/pkg/handlers"
	"social/queries"
	"social/queries/middleware"
	ws "social/queries/websocket"

	"github.com/gorilla/mux"
)

// Server represents the API server
type Server struct {
	mux *mux.Router
	db  *sql.DB
	hub *ws.Hub
}

// NewServer creates a new server instance
func NewServer(db *sql.DB) *Server {
	router := mux.NewRouter()
	hub := ws.NewHub()
	go hub.Run() // Start the WebSocket hub

	server := &Server{
		mux: router,
		db:  db,
		hub: hub,
	}
	queries.SetNotificationSender(func(userID string, notification map[string]interface{}) {
		message, err := json.Marshal(map[string]interface{}{
			"type": "notification",
			"data": notification,
		})
		if err == nil {
			hub.SendToUser(userID, message)
		}
	})

	server.setupRoutes()
	return server
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Initialize queries
	groupQueries := queries.NewGroupQueries(s.db)

	// Initialize handlers
	groupHandlers := feature.NewGroupHandlers(groupQueries)

	// Base API router
	api := s.mux.PathPrefix("/api").Subrouter()
	api.Use(handlers.CORSMiddleware)

	// Handle OPTIONS requests globally for the API subrouter
	api.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Public routes (No authentication required)
	api.HandleFunc("/register", handlers.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/login", handlers.Login).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/register", handlers.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/login", handlers.Login).Methods("POST", "OPTIONS")

	// Protected routes (Require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// User / Auth
	protected.HandleFunc("/logout", handlers.Logout).Methods("POST", "OPTIONS")
	protected.HandleFunc("/me", handlers.GetMe).Methods("GET", "OPTIONS")
	protected.HandleFunc("/session", handlers.GetSession).Methods("GET", "OPTIONS") // From auth.go
	protected.HandleFunc("/auth/logout", handlers.Logout).Methods("POST", "OPTIONS")
	protected.HandleFunc("/auth/me", handlers.GetMe).Methods("GET", "OPTIONS")
	protected.HandleFunc("/auth/session", handlers.GetSession).Methods("GET", "OPTIONS")

	// Profiles
	protected.HandleFunc("/users/{id}", handlers.GetUserProfile).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/{id}", handlers.UpdateUserProfile).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/users/{id}/privacy", handlers.UpdateProfilePrivacy).Methods("PUT", "OPTIONS")

	// Posts
	protected.HandleFunc("/posts", feature.GetFeedHandler).Methods("GET", "OPTIONS")
	protected.HandleFunc("/posts", feature.CreatePostHandler).Methods("POST", "OPTIONS")
	protected.HandleFunc("/users/{id}/posts", feature.GetUserPostsHandler).Methods("GET", "OPTIONS")
	protected.HandleFunc("/posts/{id}/comments", feature.GetCommentsHandler).Methods("GET", "OPTIONS")
	protected.HandleFunc("/posts/{id}/comments", feature.CreateCommentHandler).Methods("POST", "OPTIONS")

	// Followers
	protected.HandleFunc("/users/{id}/follow", feature.FollowUserHandler).Methods("POST", "OPTIONS")
	protected.HandleFunc("/users/{id}/follow", feature.UnfollowHandler).Methods("DELETE", "OPTIONS")
	protected.HandleFunc("/followers/{id}/accept", feature.AcceptFollowHandler).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/followers/{id}/decline", feature.DeclineFollowHandler).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/users/{id}/followers", feature.GetFollowersHandler).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/{id}/following", feature.GetFollowingHandler).Methods("GET", "OPTIONS")

	// Notifications
	protected.HandleFunc("/notifications", handlers.GetNotifications).Methods("GET", "OPTIONS")
	protected.HandleFunc("/notifications/{notificationId}/read", handlers.MarkNotificationAsRead).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/notifications/read-all", handlers.MarkAllNotificationsAsRead).Methods("PUT", "OPTIONS")

	// Chat routes
	protected.HandleFunc("/chat/users", handlers.GetDMEligibleUsers).Methods("GET", "OPTIONS")
	protected.HandleFunc("/chat/{userId}", handlers.GetPrivateMessageHistory).Methods("GET", "OPTIONS")

	// Groups routes
	protected.HandleFunc("/groups", groupHandlers.ListGroups).Methods("GET", "OPTIONS")
	protected.HandleFunc("/groups", groupHandlers.CreateGroup).Methods("POST", "OPTIONS")
	protected.HandleFunc("/groups/{id}", groupHandlers.GetGroupDetails).Methods("GET", "OPTIONS")
	protected.HandleFunc("/groups/{id}/invite", groupHandlers.InviteUserToGroup).Methods("POST", "OPTIONS")
	protected.HandleFunc("/groups/{id}/request", groupHandlers.RequestToJoinGroup).Methods("POST", "OPTIONS")
	protected.HandleFunc("/groups/{id}/events", groupHandlers.ListGroupEvents).Methods("GET", "OPTIONS")
	protected.HandleFunc("/groups/{id}/events", groupHandlers.CreateEvent).Methods("POST", "OPTIONS")

	// Group Members endpoints
	protected.HandleFunc("/groups/{id}/members/{userId}/accept", groupHandlers.AcceptMember).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/groups/{id}/members/{userId}/decline", groupHandlers.DeclineMember).Methods("PUT", "OPTIONS")

	// Group Invites endpoints
	protected.HandleFunc("/group-invites/{id}/accept", groupHandlers.AcceptGroupInvite).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/group-invites/{id}/decline", groupHandlers.DeclineGroupInvite).Methods("PUT", "OPTIONS")

	// Group Messages
	protected.HandleFunc("/groups/{groupId}/messages", handlers.GetGroupMessageHistory).Methods("GET", "OPTIONS")

	// Events routes
	protected.HandleFunc("/events/{id}/respond", groupHandlers.RespondToEvent).Methods("POST", "OPTIONS")

	// WebSockets (doesn't start with /api)
	wsRouter := s.mux.PathPrefix("/ws").Subrouter()
	wsRouter.Use(handlers.CORSMiddleware)
	wsRouter.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	wsRouter.HandleFunc("", handlers.HandleWebSocketUpgrade(s.hub)).Methods("GET", "OPTIONS")

	s.mux.PathPrefix("/uploads/").Handler(http.FileServer(http.Dir(".")))
}

// Start starts the server
func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}
