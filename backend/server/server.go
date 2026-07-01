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
)

type Server struct {
	mux *http.ServeMux
	db  *sql.DB
	hub *ws.Hub
}

func NewServer(db *sql.DB) *Server {
	mux := http.NewServeMux()
	hub := ws.NewHub()
	go hub.Run()

	server := &Server{
		mux: mux,
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

func (s *Server) setupRoutes() {
	groupQueries := queries.NewGroupQueries(s.db)
	groupHandlers := feature.NewGroupHandlers(groupQueries)

	auth := middleware.AuthMiddleware

	// Public routes
	s.mux.HandleFunc("POST /api/register", handlers.Register)
	s.mux.HandleFunc("POST /api/login", handlers.Login)
	s.mux.HandleFunc("POST /api/auth/register", handlers.Register)
	s.mux.HandleFunc("POST /api/auth/login", handlers.Login)

	// User / Auth
	s.mux.Handle("POST /api/logout", auth(http.HandlerFunc(handlers.Logout)))
	s.mux.Handle("GET /api/me", auth(http.HandlerFunc(handlers.GetMe)))
	s.mux.Handle("GET /api/session", auth(http.HandlerFunc(handlers.GetSession)))
	s.mux.Handle("POST /api/auth/logout", auth(http.HandlerFunc(handlers.Logout)))
	s.mux.Handle("GET /api/auth/me", auth(http.HandlerFunc(handlers.GetMe)))
	s.mux.Handle("GET /api/auth/session", auth(http.HandlerFunc(handlers.GetSession)))

	// Profiles
	s.mux.Handle("GET /api/users/{id}", auth(http.HandlerFunc(handlers.GetUserProfile)))
	s.mux.Handle("PUT /api/users/{id}", auth(http.HandlerFunc(handlers.UpdateUserProfile)))
	s.mux.Handle("PUT /api/users/{id}/privacy", auth(http.HandlerFunc(handlers.UpdateProfilePrivacy)))

	// Posts
	s.mux.Handle("GET /api/posts", auth(http.HandlerFunc(feature.GetFeedHandler)))
	s.mux.Handle("POST /api/posts", auth(http.HandlerFunc(feature.CreatePostHandler)))
	s.mux.Handle("GET /api/users/{id}/posts", auth(http.HandlerFunc(feature.GetUserPostsHandler)))
	s.mux.Handle("GET /api/posts/{id}/comments", auth(http.HandlerFunc(feature.GetCommentsHandler)))
	s.mux.Handle("POST /api/posts/{id}/comments", auth(http.HandlerFunc(feature.CreateCommentHandler)))
	s.mux.Handle("PUT /api/posts/{id}", auth(http.HandlerFunc(feature.UpdatePostHandler)))
	s.mux.Handle("DELETE /api/posts/{id}", auth(http.HandlerFunc(feature.DeletePostHandler)))
	s.mux.Handle("POST /api/posts/{id}/react", auth(http.HandlerFunc(handlers.ReactToPost)))

	// Followers
	s.mux.Handle("POST /api/users/{id}/follow", auth(http.HandlerFunc(feature.FollowUserHandler)))
	s.mux.Handle("DELETE /api/users/{id}/follow", auth(http.HandlerFunc(feature.UnfollowHandler)))
	s.mux.Handle("PUT /api/followers/{id}/accept", auth(http.HandlerFunc(feature.AcceptFollowHandler)))
	s.mux.Handle("PUT /api/followers/{id}/decline", auth(http.HandlerFunc(feature.DeclineFollowHandler)))
	s.mux.Handle("GET /api/users/{id}/followers", auth(http.HandlerFunc(feature.GetFollowersHandler)))
	s.mux.Handle("GET /api/users/{id}/following", auth(http.HandlerFunc(feature.GetFollowingHandler)))

	// Notifications
	s.mux.Handle("GET /api/notifications", auth(http.HandlerFunc(handlers.GetNotifications)))
	s.mux.Handle("PUT /api/notifications/{notificationId}/read", auth(http.HandlerFunc(handlers.MarkNotificationAsRead)))
	s.mux.Handle("PUT /api/notifications/read-all", auth(http.HandlerFunc(handlers.MarkAllNotificationsAsRead)))

	// Chat routes
	s.mux.Handle("GET /api/chat/users", auth(http.HandlerFunc(handlers.GetDMEligibleUsers)))
	s.mux.Handle("GET /api/chat/{userId}", auth(http.HandlerFunc(handlers.GetPrivateMessageHistory)))
	s.mux.Handle("POST /api/chat/{userId}/messages", auth(handlers.SendPrivateMessage(s.hub)))

	// Groups routes
	s.mux.Handle("GET /api/groups", auth(http.HandlerFunc(groupHandlers.ListGroups)))
	s.mux.Handle("POST /api/groups", auth(http.HandlerFunc(groupHandlers.CreateGroup)))
	s.mux.Handle("GET /api/groups/{id}", auth(http.HandlerFunc(groupHandlers.GetGroupDetails)))
	s.mux.Handle("POST /api/groups/{id}/invite", auth(http.HandlerFunc(groupHandlers.InviteUserToGroup)))
	s.mux.Handle("POST /api/groups/{id}/request", auth(http.HandlerFunc(groupHandlers.RequestToJoinGroup)))
	s.mux.Handle("GET /api/groups/{id}/posts", auth(http.HandlerFunc(feature.GetGroupPostsHandler)))
	s.mux.Handle("GET /api/groups/{id}/events", auth(http.HandlerFunc(groupHandlers.ListGroupEvents)))
	s.mux.Handle("POST /api/groups/{id}/events", auth(http.HandlerFunc(groupHandlers.CreateEvent)))
	s.mux.Handle("PUT /api/groups/{id}/members/{userId}/accept", auth(http.HandlerFunc(groupHandlers.AcceptMember)))
	s.mux.Handle("PUT /api/groups/{id}/members/{userId}/decline", auth(http.HandlerFunc(groupHandlers.DeclineMember)))
	s.mux.Handle("PUT /api/group-invites/{id}/accept", auth(http.HandlerFunc(groupHandlers.AcceptGroupInvite)))
	s.mux.Handle("PUT /api/group-invites/{id}/decline", auth(http.HandlerFunc(groupHandlers.DeclineGroupInvite)))

	// Group Messages
	s.mux.Handle("GET /api/groups/{groupId}/messages", auth(http.HandlerFunc(handlers.GetGroupMessageHistory)))
	s.mux.Handle("POST /api/groups/{groupId}/messages", auth(handlers.SendGroupMessage(s.hub)))

	// Events routes
	s.mux.Handle("POST /api/events/{id}/respond", auth(http.HandlerFunc(groupHandlers.RespondToEvent)))

	// WebSocket (no auth middleware - it does its own session validation)
	s.mux.HandleFunc("GET /ws", handlers.HandleWebSocketUpgrade(s.hub))

	// Uploads
	s.mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, handlers.CORSMiddleware(s.mux))
}
