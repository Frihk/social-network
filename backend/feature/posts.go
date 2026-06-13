package feature

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"social/models"
	"social/queries"
	"social/queries/utils"
)

// GetFeedHandler handles GET /api/posts
func GetFeedHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		posts, err := queries.GetFeed(db, userID)
		if err != nil {
			http.Error(w, "failed to fetch feed", http.StatusInternalServerError)
			return
		}

		if posts == nil {
			posts = []models.Post{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

// CreatePostHandler handles POST /api/posts
func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			return
		}

		content := strings.TrimSpace(r.FormValue("content"))
		if content == "" {
			http.Error(w, "content is required", http.StatusBadRequest)
			return
		}

		privacy := r.FormValue("privacy")
		if privacy != "public" && privacy != "almost_private" && privacy != "private" {
			http.Error(w, "invalid privacy value", http.StatusBadRequest)
			return
		}

		post := models.Post{
			UserID:  userID,
			Content: content,
			Privacy: privacy,
		}

		// handle optional image
		file, header, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath, err := utils.SaveImage(file, header)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			post.ImagePath = &imagePath
		}

		// parse allowed viewers if privacy is private
		var allowedViewers []int
		if privacy == "private" {
			viewerStrs := r.Form["allowed_viewers"]
			for _, v := range viewerStrs {
				id, err := strconv.Atoi(v)
				if err != nil {
					http.Error(w, "invalid allowed viewer ID", http.StatusBadRequest)
					return
				}
				allowedViewers = append(allowedViewers, id)
			}
		}

		postID, err := queries.CreatePost(db, post, allowedViewers)
		if err != nil {
			http.Error(w, "failed to create post", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int64{"id": postID})
	}
}

// GetUserPostsHandler handles GET /api/users/:id/posts
func GetUserPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// extract :id from /api/users/:id/posts
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "invalid URL", http.StatusBadRequest)
			return
		}
		targetUserID, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		posts, err := queries.GetPostsByUserID(db, targetUserID, currentUserID)
		if err != nil {
			http.Error(w, "failed to fetch posts", http.StatusInternalServerError)
			return
		}

		if posts == nil {
			posts = []models.Post{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

// CreateCommentHandler handles POST /api/posts/:id/comments
func CreateCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// extract :id from /api/posts/:id/comments
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "invalid URL", http.StatusBadRequest)
			return
		}
		postID, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "invalid post ID", http.StatusBadRequest)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			return
		}

		content := strings.TrimSpace(r.FormValue("content"))
		if content == "" {
			http.Error(w, "content is required", http.StatusBadRequest)
			return
		}

		comment := models.Comment{
			PostID:  postID,
			UserID:  userID,
			Content: content,
		}

		// handle optional image
		file, header, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imagePath, err := utils.SaveImage(file, header)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			comment.ImagePath = &imagePath
		}

		commentID, err := queries.CreateComment(db, comment)
		if err != nil {
			http.Error(w, "failed to create comment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int64{"id": commentID})
	}
}

// GetCommentsHandler handles GET /api/posts/:id/comments
func GetCommentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// extract :id from /api/posts/:id/comments
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "invalid URL", http.StatusBadRequest)
			return
		}
		postID, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "invalid post ID", http.StatusBadRequest)
			return
		}

		comments, err := queries.GetCommentsByPostID(db, postID)
		if err != nil {
			http.Error(w, "failed to fetch comments", http.StatusInternalServerError)
			return
		}

		if comments == nil {
			comments = []models.Comment{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}
