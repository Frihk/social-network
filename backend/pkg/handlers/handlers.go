package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

// UserIDKey is used to store user ID in request context
type UserIDKey struct{}

// AuthMiddleware ensures user_id is present in the request
// This is a basic middleware that checks for user_id header
// In production, this should validate JWT tokens or sessions
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			// Try to get from query parameter (for development)
			userID = r.URL.Query().Get("user_id")
		}

		if userID == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "user_id required"})
			return
		}

		// Store in context for use in handlers
		ctx := context.WithValue(r.Context(), UserIDKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORSMiddleware enables CORS
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-User-ID")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// JSONResponse writes a JSON response
func JSONResponse(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// JSONError writes a JSON error response
func JSONError(w http.ResponseWriter, message string, statusCode int) error {
	return JSONResponse(w, map[string]string{"error": message}, statusCode)
}
