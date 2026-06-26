package handlers

import (
	"encoding/json"
	"net/http"

	"social/queries"
	"social/queries/middleware"
)

// GetUserGroups returns all groups the logged-in user is a member of
func GetUserGroups(w http.ResponseWriter, r *http.Request) {
	loggedInUserID := middleware.GetUserID(r.Context())
	if loggedInUserID == "" {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	groups, err := queries.GetUserGroups(loggedInUserID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch groups"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groups)
}
