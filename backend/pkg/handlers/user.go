package handlers

import (
	"encoding/json"
	"net/http"

	"social/models"
	"social/pkg/db/sqlite"
	"social/queries/middleware"

	"github.com/gorilla/mux"
)

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := models.GetUserByID(sqlite.DB, userID)
	if err != nil {
		http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
		return
	}

	currentUserID := middleware.GetUserID(r.Context())

	if user.IsPrivate && currentUserID != user.ID {
		// TODO: Check if current user is a follower
		// For now, return limited info
		limitedUser := map[string]interface{}{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"avatar":     user.Avatar,
			"is_private": user.IsPrivate,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(limitedUser)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateProfilePrivacy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileUserID := vars["id"]

	currentUserID := middleware.GetUserID(r.Context())
	if currentUserID != profileUserID {
		http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
		return
	}

	var req struct {
		IsPrivate bool `json:"is_private"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := models.SetProfilePrivacy(sqlite.DB, currentUserID, req.IsPrivate); err != nil {
		http.Error(w, `{"error":"Failed to update privacy"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"is_private": req.IsPrivate})
}
