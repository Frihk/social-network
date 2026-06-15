package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"context"
	"fmt"
	// "path/filepath"
	"os"
	// "errors"

	"social/feature"
	_ "github.com/mattn/go-sqlite3"
	// "github.com/golang-migrate/migrate/v4"
	// "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// temporary auth middleware for testing — injects a hardcoded user_id
func tempAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// hardcoded user ID 1 for testing
		ctx := context.WithValue(r.Context(), "user_id", 1)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	// run migrations by reading SQL files directly
	migrations := []string{
		"pkg/db/migration/sqlite/000001_create_users_table.up.sql",
		"pkg/db/migration/sqlite/000002_create_posts_table.up.sql",
		"pkg/db/migration/sqlite/000004_create_followers_table.up.sql",
	}

	for _, f := range migrations {
		content, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", f, err)
		}
		if _, err := db.Exec(string(content)); err != nil {
			log.Printf("skipping %s: %v", f, err)
		}
	}

	log.Println("database ready")
	return db, nil
}
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func main() {
	db, err := initDB("./test.db")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	// post routes
	mux.Handle("/api/posts", tempAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			feature.GetFeedHandler(db)(w, r)
		case http.MethodPost:
			feature.CreatePostHandler(db)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/api/posts/", tempAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			feature.GetCommentsHandler(db)(w, r)
		case http.MethodPost:
			feature.CreateCommentHandler(db)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// user posts and follow routes
	mux.Handle("/api/users/", tempAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case len(path) > 0 && path[len(path)-6:] == "/posts":
			feature.GetUserPostsHandler(db)(w, r)
		case len(path) > 0 && path[len(path)-7:] == "/follow":
			feature.FollowUserHandler(db)(w, r)
		case len(path) > 0 && path[len(path)-9:] == "/unfollow":
			feature.UnfollowHandler(db)(w, r)
		case len(path) > 0 && path[len(path)-10:] == "/followers":
			feature.GetFollowersHandler(db)(w, r)
		case len(path) > 0 && path[len(path)-10:] == "/following":
			feature.GetFollowingHandler(db)(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	})))

	// follower accept/decline routes
	mux.Handle("/api/followers/", tempAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case len(path) > 0 && path[len(path)-7:] == "/accept":
			feature.AcceptFollowHandler(db)(w, r)
		case len(path) > 0 && path[len(path)-8:] == "/decline":
			feature.DeclineFollowHandler(db)(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	})))

	log.Println("test server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}