package server

import (
	"database/sql"
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
	api.HandleFunc("/register", handlers.Register).Methods("POST")
	api.HandleFunc("/login", handlers.Login).Methods("POST")
	
	// Protected routes (Require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	
	// User / Auth
	protected.HandleFunc("/logout", handlers.Logout).Methods("POST")
	protected.HandleFunc("/me", handlers.GetMe).Methods("GET")
	protected.HandleFunc("/session", handlers.GetSession).Methods("GET") // From auth.go

	// Chat routes
	protected.HandleFunc("/chat/users", handlers.GetDMEligibleUsers).Methods("GET")
	protected.HandleFunc("/chat/{userId}", handlers.GetPrivateMessageHistory).Methods("GET")

	// Groups routes
	protected.HandleFunc("/groups", groupHandlers.ListGroups).Methods("GET")
	protected.HandleFunc("/groups", groupHandlers.CreateGroup).Methods("POST")
	protected.HandleFunc("/groups/{id}", groupHandlers.GetGroupDetails).Methods("GET")
	protected.HandleFunc("/groups/{id}/invite", groupHandlers.InviteUserToGroup).Methods("POST")
	protected.HandleFunc("/groups/{id}/request", groupHandlers.RequestToJoinGroup).Methods("POST")
	protected.HandleFunc("/groups/{id}/events", groupHandlers.ListGroupEvents).Methods("GET")
	protected.HandleFunc("/groups/{id}/events", groupHandlers.CreateEvent).Methods("POST")
	
	// Group Members endpoints
	protected.HandleFunc("/groups/{id}/members/{userId}/accept", groupHandlers.AcceptMember).Methods("PUT")
	protected.HandleFunc("/groups/{id}/members/{userId}/decline", groupHandlers.DeclineMember).Methods("PUT")
	
	// Group Invites endpoints
	protected.HandleFunc("/group-invites/{id}/accept", groupHandlers.AcceptGroupInvite).Methods("PUT")
	protected.HandleFunc("/group-invites/{id}/decline", groupHandlers.DeclineGroupInvite).Methods("PUT")
	
	// Group Messages
	protected.HandleFunc("/groups/{groupId}/messages", handlers.GetGroupMessageHistory).Methods("GET")

	// Events routes
	protected.HandleFunc("/events/{id}/respond", groupHandlers.RespondToEvent).Methods("POST")

	// WebSockets (doesn't start with /api)
	wsRouter := s.mux.PathPrefix("/ws").Subrouter()
	wsRouter.Use(handlers.CORSMiddleware)
	wsRouter.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	wsRouter.HandleFunc("", handlers.HandleWebSocketUpgrade(s.hub)).Methods("GET")
}

// Start starts the server
func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}
