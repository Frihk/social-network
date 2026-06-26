package server

import (
	"log"
	"net/http"
	"social/feature"
	"social/pkg/handlers"
	"social/queries/middleware"

	"github.com/gorilla/mux"
)

func NewServer() *http.Server {
	router := mux.NewRouter()

	// Apply CORS middleware globally
	router.Use(middleware.CORSMiddleware)

	// --- Public routes ---
	router.HandleFunc("/api/auth/register", handlers.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/login", handlers.Login).Methods("POST", "OPTIONS")

	// --- Protected routes ---
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Person 1 (auth + users)
	protected.HandleFunc("/auth/logout", handlers.Logout).Methods("POST", "OPTIONS")
	protected.HandleFunc("/auth/me", handlers.GetMe).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/{id}", handlers.GetUserProfile).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/{id}/privacy", handlers.UpdateProfilePrivacy).Methods("PUT", "OPTIONS")

	// Person 2 (posts + comments)
	protected.HandleFunc("/posts", feature.GetFeedHandler).Methods("GET", "OPTIONS")
	protected.HandleFunc("/posts", feature.CreatePostHandler).Methods("POST", "OPTIONS")
	protected.HandleFunc("/users/{id}/posts", feature.GetUserPostsHandler).Methods("GET", "OPTIONS")
	protected.HandleFunc("/posts/{id}/comments", feature.CreateCommentHandler).Methods("POST", "OPTIONS")
	protected.HandleFunc("/posts/{id}/comments", feature.GetCommentsHandler).Methods("GET", "OPTIONS")

	// Person 2 (followers)
	protected.HandleFunc("/users/{id}/follow", feature.FollowUserHandler).Methods("POST", "OPTIONS")
	protected.HandleFunc("/users/{id}/follow", feature.UnfollowHandler).Methods("DELETE", "OPTIONS")
	protected.HandleFunc("/followers/{id}/accept", feature.AcceptFollowHandler).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/followers/{id}/decline", feature.DeclineFollowHandler).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/users/{id}/followers", feature.GetFollowersHandler).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/{id}/following", feature.GetFollowingHandler).Methods("GET", "OPTIONS")

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server configured on :8080")
	return server
}
