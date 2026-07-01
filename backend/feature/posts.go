package feature

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"social/models"
	"social/pkg/db/sqlite"
	"social/queries"
	"social/queries/middleware"
	"social/queries/utils"
)

func GetFeedHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	posts, err := queries.GetFeed(sqlite.DB, userID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch feed"}`, http.StatusInternalServerError)
		return
	}

	if posts == nil {
		posts = []models.Post{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error":"failed to parse form"}`, http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Error(w, `{"error":"content is required"}`, http.StatusBadRequest)
		return
	}

	privacy := r.FormValue("privacy")
	if privacy != "public" && privacy != "almost_private" && privacy != "private" {
		http.Error(w, `{"error":"invalid privacy value"}`, http.StatusBadRequest)
		return
	}

	groupID := r.FormValue("group_id")

	post := models.Post{
		UserID:  userID,
		Content: content,
		Privacy: privacy,
	}
	if groupID != "" {
		post.GroupID = &groupID
	}

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

	var allowedViewers []string
	if privacy == "private" {
		allowedViewers = r.Form["allowed_viewers"]
	}

	postID, err := queries.CreatePost(sqlite.DB, post, allowedViewers)
	if err != nil {
		http.Error(w, `{"error":"failed to create post"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": postID})
}

func GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID := middleware.GetUserID(r.Context())
	if currentUserID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	targetUserID := r.PathValue("id")
	if targetUserID == "" {
		http.Error(w, `{"error":"invalid user ID"}`, http.StatusBadRequest)
		return
	}

	posts, err := queries.GetPostsByUserID(sqlite.DB, targetUserID, currentUserID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch posts"}`, http.StatusInternalServerError)
		return
	}

	if posts == nil {
		posts = []models.Post{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	postIDStr := r.PathValue("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid post ID"}`, http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error":"failed to parse form"}`, http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Error(w, `{"error":"content is required"}`, http.StatusBadRequest)
		return
	}

	privacy := r.FormValue("privacy")
	if privacy != "public" && privacy != "almost_private" && privacy != "private" {
		http.Error(w, `{"error":"invalid privacy value"}`, http.StatusBadRequest)
		return
	}

	if err := queries.UpdatePost(sqlite.DB, postID, content, privacy, userID); err != nil {
		http.Error(w, `{"error":"failed to update post"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	postIDStr := r.PathValue("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid post ID"}`, http.StatusBadRequest)
		return
	}

	if err := queries.DeletePost(sqlite.DB, postID, userID); err != nil {
		http.Error(w, `{"error":"failed to delete post"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func GetGroupPostsHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	groupID := r.PathValue("id")
	if groupID == "" {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	posts, err := queries.GetPostsByGroupID(sqlite.DB, groupID, userID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch group posts"}`, http.StatusInternalServerError)
		return
	}

	if posts == nil {
		posts = []models.Post{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error":"invalid post ID"}`, http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error":"failed to parse form"}`, http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Error(w, `{"error":"content is required"}`, http.StatusBadRequest)
		return
	}

	comment := models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}

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

	commentID, err := queries.CreateComment(sqlite.DB, comment)
	if err != nil {
		http.Error(w, `{"error":"failed to create comment"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": commentID})
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error":"invalid post ID"}`, http.StatusBadRequest)
		return
	}

	comments, err := queries.GetCommentsByPostID(sqlite.DB, postID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch comments"}`, http.StatusInternalServerError)
		return
	}

	if comments == nil {
		comments = []models.Comment{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}