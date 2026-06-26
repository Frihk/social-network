package handlers

import (
	"encoding/json"
	"net/http"

	"social/queries"
	"social/queries/middleware"

	"github.com/gorilla/mux"
)

// GetNotifications returns all notifications for the logged-in user
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	loggedInUserID := middleware.GetUserID(r.Context())
	if loggedInUserID == "" {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	notifications, err := queries.GetNotificationsForUser(loggedInUserID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch notifications"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   notifications,
	})
}

// MarkNotificationAsRead marks a single notification as read
// Only allows the owner of the notification to mark it as read
func MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notificationID := vars["notificationId"]

	loggedInUserID := middleware.GetUserID(r.Context())
	if loggedInUserID == "" {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if notificationID == "" {
		http.Error(w, `{"error":"Notification ID is required"}`, http.StatusBadRequest)
		return
	}

	// Verify that the notification belongs to the logged-in user
	ownerID, err := queries.GetNotificationOwner(notificationID)
	if err != nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusForbidden)
		return
	}

	if ownerID != loggedInUserID {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusForbidden)
		return
	}

	if err := queries.MarkAsRead(notificationID); err != nil {
		http.Error(w, `{"error":"Failed to mark notification as read"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Notification marked as read",
	})
}

// MarkAllNotificationsAsRead marks all notifications for the logged-in user as read
func MarkAllNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	loggedInUserID := middleware.GetUserID(r.Context())
	if loggedInUserID == "" {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if err := queries.MarkAllAsRead(loggedInUserID); err != nil {
		http.Error(w, `{"error":"Failed to mark all notifications as read"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "All notifications marked as read",
	})
}