package server

import (
	"log"
	"net/http"

	"social/pkg/handlers"
	"social/queries/middleware"

	"github.com/gorilla/mux"
)

func NewServer() *http.Server {
	router := mux.NewRouter()

	// Apply CORS middleware globally
	router.Use(middleware.CORSMiddleware)

	// Public routes
	router.HandleFunc("/api/auth/register", handlers.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/login", handlers.Login).Methods("POST", "OPTIONS")

	// Protected routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/auth/logout", handlers.Logout).Methods("POST", "OPTIONS")
	protected.HandleFunc("/auth/me", handlers.GetMe).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/{id}", handlers.GetUserProfile).Methods("GET", "OPTIONS")
	protected.HandleFunc("/users/{id}/privacy", handlers.UpdateProfilePrivacy).Methods("PUT", "OPTIONS")

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server configured on :8080")
	return server
}
