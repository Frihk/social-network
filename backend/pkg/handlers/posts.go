package handlers

import (
	"encoding/json"
	"net/http"
	"social/models"
	"social/pkg/db/sqlite"
	"social/queries"
	"strings"
)

// ReactToPost handles adding or removing an emoji reaction to a post
func ReactToPost(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	postID := pathParts[len(pathParts)-2] // e.g., /api/posts/{id}/react

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := queries.ToggleReaction(sqlite.DB, postID, userID, req.Emoji)
	if err != nil {
		http.Error(w, "Failed to process reaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reaction updated successfully"})
}
