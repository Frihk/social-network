package handlers

import (
	"encoding/json"
	"net/http"

	"social/queries"
	"social/queries/middleware"

	"github.com/gorilla/mux"
)

// GetDMEligibleUsers returns all users the logged-in user can DM
// (users where at least one follows the other with accepted status)
func GetDMEligibleUsers(w http.ResponseWriter, r *http.Request) {
	loggedInUserID := middleware.GetUserID(r.Context())
	if loggedInUserID == "" {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	users, err := queries.GetDMEligibleUsers(loggedInUserID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch DM eligible users"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetPrivateMessageHistory returns the message history between two users
// Only returns messages if there is an accepted follow relationship
func GetPrivateMessageHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	loggedInUserID := middleware.GetUserID(r.Context())
	if loggedInUserID == "" {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if userID == "" {
		http.Error(w, `{"error":"User ID is required"}`, http.StatusBadRequest)
		return
	}

	// Check if there is a follow relationship
	if !queries.CanDM(loggedInUserID, userID) {
		http.Error(w, `{"error":"No follow relationship with this user"}`, http.StatusForbidden)
		return
	}

	messages, err := queries.GetPrivateMessageHistory(loggedInUserID, userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch message history"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

// GetGroupMessageHistory returns the message history for a group
// Only returns messages if the user is an accepted member of the group
func GetGroupMessageHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupId"]

	loggedInUserID := middleware.GetUserID(r.Context())
	if loggedInUserID == "" {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if groupID == "" {
		http.Error(w, `{"error":"Group ID is required"}`, http.StatusBadRequest)
		return
	}

	// Check if user is a member of the group
	if !queries.IsGroupMember(loggedInUserID, groupID) {
		http.Error(w, `{"error":"Not a member of this group"}`, http.StatusForbidden)
		return
	}

	messages, err := queries.GetGroupMessageHistory(groupID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch group message history"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}