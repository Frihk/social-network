package feature

import (
	"encoding/json"
	"net/http"
	"strconv"

	"social/models"
	"social/pkg/db/sqlite"
	"social/queries"
	"social/queries/middleware"
)

func FollowUserHandler(w http.ResponseWriter, r *http.Request) {
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

	if currentUserID == targetUserID {
		http.Error(w, `{"error":"cannot follow yourself"}`, http.StatusBadRequest)
		return
	}

	existing, err := queries.GetFollowStatus(sqlite.DB, currentUserID, targetUserID)
	if err == nil {
		if existing.Status == "accepted" {
			http.Error(w, `{"error":"already following"}`, http.StatusConflict)
			return
		}
		if existing.Status == "pending" {
			http.Error(w, `{"error":"follow request already pending"}`, http.StatusConflict)
			return
		}
	}

	targetUser, err := models.GetUserByID(sqlite.DB, targetUserID)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	status := "accepted"
	if targetUser.IsPrivate {
		status = "pending"
	}

	followerRowID, err := queries.CreateFollower(sqlite.DB, currentUserID, targetUserID, status)
	if err != nil {
		http.Error(w, `{"error":"failed to follow user"}`, http.StatusInternalServerError)
		return
	}

	if status == "pending" {
		if err := queries.CreateNotification(targetUserID, "follow_request", currentUserID, strconv.FormatInt(followerRowID, 10)); err != nil {
			http.Error(w, `{"error":"failed to create follow request notification"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func AcceptFollowHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID := middleware.GetUserID(r.Context())
	if currentUserID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	followerRowID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error":"invalid follower ID"}`, http.StatusBadRequest)
		return
	}

	follower, err := queries.GetFollowerByID(sqlite.DB, followerRowID)
	if err != nil {
		http.Error(w, `{"error":"follow request not found"}`, http.StatusNotFound)
		return
	}
	if follower.FollowingID != currentUserID {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}

	if err := queries.AcceptFollowRequest(sqlite.DB, followerRowID); err != nil {
		http.Error(w, `{"error":"failed to accept follow request"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}

func DeclineFollowHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID := middleware.GetUserID(r.Context())
	if currentUserID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	followerRowID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, `{"error":"invalid follower ID"}`, http.StatusBadRequest)
		return
	}

	follower, err := queries.GetFollowerByID(sqlite.DB, followerRowID)
	if err != nil {
		http.Error(w, `{"error":"follow request not found"}`, http.StatusNotFound)
		return
	}
	if follower.FollowingID != currentUserID {
		http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		return
	}

	if err := queries.DeclineFollowRequest(sqlite.DB, followerRowID); err != nil {
		http.Error(w, `{"error":"failed to decline follow request"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "declined"})
}

func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := queries.Unfollow(sqlite.DB, currentUserID, targetUserID); err != nil {
		http.Error(w, `{"error":"failed to unfollow user"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "unfollowed"})
}

func GetFollowersHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, `{"error":"invalid user ID"}`, http.StatusBadRequest)
		return
	}

	followers, err := queries.GetFollowers(sqlite.DB, userID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch followers"}`, http.StatusInternalServerError)
		return
	}

	if followers == nil {
		followers = []models.Follower{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func GetFollowingHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, `{"error":"invalid user ID"}`, http.StatusBadRequest)
		return
	}

	following, err := queries.GetFollowing(sqlite.DB, userID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch following"}`, http.StatusInternalServerError)
		return
	}

	if following == nil {
		following = []models.Follower{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}
