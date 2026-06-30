package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"social/queries"
	"social/queries/middleware"
	ws "social/queries/websocket"
)

type sendMessageRequest struct {
	Content string `json:"content"`
}

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
	userID := r.PathValue("userId")

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

// SendPrivateMessage saves and broadcasts a direct message.
func SendPrivateMessage(hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("userId")

		loggedInUserID := middleware.GetUserID(r.Context())
		if loggedInUserID == "" {
			http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		if userID == "" {
			http.Error(w, `{"error":"User ID is required"}`, http.StatusBadRequest)
			return
		}

		if !queries.CanDM(loggedInUserID, userID) {
			http.Error(w, `{"error":"No follow relationship with this user"}`, http.StatusForbidden)
			return
		}

		var req sendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		if req.Content == "" {
			http.Error(w, `{"error":"Message content is required"}`, http.StatusBadRequest)
			return
		}

		messageID, err := queries.SavePrivateMessage(loggedInUserID, userID, req.Content)
		if err != nil {
			http.Error(w, `{"error":"Failed to save message"}`, http.StatusInternalServerError)
			return
		}

		message := ws.Message{
			Type:       "private_message",
			ID:         messageID,
			UserID:     userID,
			SenderID:   loggedInUserID,
			ReceiverID: userID,
			Content:    req.Content,
			CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		}

		if data, err := json.Marshal(message); err == nil {
			hub.SendToUser(loggedInUserID, data)
			hub.SendToUser(userID, data)
		}
		_ = queries.CreateNotification(userID, "private_message", loggedInUserID, messageID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
	}
}

// GetGroupMessageHistory returns the message history for a group
// Only returns messages if the user is an accepted member of the group
func GetGroupMessageHistory(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")

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

// SendGroupMessage saves and broadcasts a group message.
func SendGroupMessage(hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupID := r.PathValue("groupId")

		loggedInUserID := middleware.GetUserID(r.Context())
		if loggedInUserID == "" {
			http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		if groupID == "" {
			http.Error(w, `{"error":"Group ID is required"}`, http.StatusBadRequest)
			return
		}

		if !queries.IsGroupMember(loggedInUserID, groupID) {
			http.Error(w, `{"error":"Not a member of this group"}`, http.StatusForbidden)
			return
		}

		var req sendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		if req.Content == "" {
			http.Error(w, `{"error":"Message content is required"}`, http.StatusBadRequest)
			return
		}

		messageID, err := queries.SaveGroupMessage(groupID, loggedInUserID, req.Content)
		if err != nil {
			http.Error(w, `{"error":"Failed to save message"}`, http.StatusInternalServerError)
			return
		}

		message := ws.Message{
			Type:      "group_message",
			ID:        messageID,
			GroupID:   groupID,
			SenderID:  loggedInUserID,
			Content:   req.Content,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		}

		if data, err := json.Marshal(message); err == nil {
			hub.BroadcastToGroup(groupID, data)
		}
		memberIDs, err := queries.GetAcceptedGroupMemberIDs(groupID)
		if err == nil {
			for _, memberID := range memberIDs {
				if memberID != loggedInUserID {
					_ = queries.CreateNotification(memberID, "group_message", loggedInUserID, groupID)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
	}
}
