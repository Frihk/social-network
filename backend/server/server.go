package server

import (
	"database/sql"
	"net/http"
	"social/feature"
	"social/pkg/handlers"
	"social/queries"
)

// Server represents the API server
type Server struct {
	mux *http.ServeMux
	db  *sql.DB
}

// NewServer creates a new server instance
func NewServer(db *sql.DB) *Server {
	mux := http.NewServeMux()
	server := &Server{
		mux: mux,
		db:  db,
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

	// Create a chain of middleware
	api := handlers.CORSMiddleware(handlers.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Route the request to the appropriate handler
		path := r.URL.Path
		method := r.Method

		// Groups endpoints
		if path == "/api/groups" {
			if method == "GET" {
				groupHandlers.ListGroups(w, r)
			} else if method == "POST" {
				groupHandlers.CreateGroup(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Get group details - /api/groups/{id}
		if len(path) > 11 && path[:11] == "/api/groups/" {
			remainder := path[11:]
			// Check if it's an ID without sub-routes
			if len(remainder) > 0 && remainder[len(remainder)-1] != '/' && !containsSlash(remainder) {
				if method == "GET" {
					groupHandlers.GetGroupDetails(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}

			// Invite user - POST /api/groups/{id}/invite
			if len(remainder) > 7 && remainder[len(remainder)-7:] == "/invite" {
				if method == "POST" {
					groupHandlers.InviteUserToGroup(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}

			// Request to join - POST /api/groups/{id}/request
			if len(remainder) > 8 && remainder[len(remainder)-8:] == "/request" {
				if method == "POST" {
					groupHandlers.RequestToJoinGroup(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}

			// Events - /api/groups/{id}/events
			if len(remainder) > 7 && remainder[len(remainder)-7:] == "/events" {
			if method == "GET" {
				groupHandlers.ListGroupEvents(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}

			// Members endpoints - /api/groups/{id}/members/{userId}/{action}
			if len(remainder) > 9 && remainder[len(remainder)-7:] == "/accept" {
				if method == "PUT" {
					groupHandlers.AcceptMember(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}

			if len(remainder) > 9 && remainder[len(remainder)-8:] == "/decline" {
				if method == "PUT" {
					groupHandlers.DeclineMember(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}
		}

		// Group invites endpoints
		if len(path) > 15 && path[:15] == "/api/group-invites/" {
			remainder := path[15:]

			// Accept invite - PUT /api/group-invites/{id}/accept
			if len(remainder) > 7 && remainder[len(remainder)-7:] == "/accept" {
				if method == "PUT" {
					groupHandlers.AcceptGroupInvite(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}

			// Decline invite - PUT /api/group-invites/{id}/decline
			if len(remainder) > 8 && remainder[len(remainder)-8:] == "/decline" {
				if method == "PUT" {
					groupHandlers.DeclineGroupInvite(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}
		}

		// Event respond - POST /api/events/{id}/respond
		if len(path) > 11 && path[:11] == "/api/events/" {
			remainder := path[11:]
			if len(remainder) > 8 && remainder[len(remainder)-8:] == "/respond" {
				if method == "POST" {
					groupHandlers.RespondToEvent(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}
		}

		http.Error(w, "Not found", http.StatusNotFound)
	})))

	s.mux.Handle("/", api)
}

// Start starts the server
func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

// containsSlash checks if a string contains a forward slash
func containsSlash(s string) bool {
	for _, c := range s {
		if c == '/' {
			return true
		}
	}
	return false
}
