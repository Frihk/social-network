package feature

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"social/queries"
	"social/models"
)

// FollowUserHandler handles POST /api/users/:id/follow
func FollowUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// extract :id from /api/users/:id/follow
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

		// prevent following yourself
		if currentUserID == targetUserID {
			http.Error(w, "cannot follow yourself", http.StatusBadRequest)
			return
		}

		// check if already following or pending
		existing, err := queries.GetFollowStatus(db, currentUserID, targetUserID)
		if err == nil {
			if existing.Status == "accepted" {
				http.Error(w, "already following", http.StatusConflict)
				return
			}
			if existing.Status == "pending" {
				http.Error(w, "follow request already pending", http.StatusConflict)
				return
			}
		}

		// check target user's profile privacy
		// TODO: replace with actual call once Person 1 pushes queries/user.go
		// targetUser, err := queries.GetUserByID(db, targetUserID)
		// if err != nil {
		// 	http.Error(w, "user not found", http.StatusNotFound)
		// 	return
		// }
		// isPrivate := targetUser.IsPrivate

		// temporary stub — assume private until Person 1's code is ready
		isPrivate := true

		status := "accepted"
		if isPrivate {
			status = "pending"
		}

		followerID, err := queries.CreateFollower(db, currentUserID, targetUserID, status)
		if err != nil {
			http.Error(w, "failed to follow user", http.StatusInternalServerError)
			return
		}

		// send notification if request is pending
		if status == "pending" {
			// TODO: uncomment once Person 4 pushes queries/notification.go
			// err = queries.CreateNotification(db, models.Notification{
			// 	UserID:   targetUserID,
			// 	Type:     "follow_request",
			// 	ActorID:  currentUserID,
			// 	EntityID: int(followerID),
			// })
			// if err != nil {
			// 	log.Printf("failed to create notification: %v", err)
			// }
			_ = followerID
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": status})
	}
}

// AcceptFollowHandler handles PUT /api/followers/:id/accept
func AcceptFollowHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// extract :id from /api/followers/:id/accept
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "invalid URL", http.StatusBadRequest)
			return
		}
		followerRowID, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "invalid follower ID", http.StatusBadRequest)
			return
		}

		// verify the current user is the recipient of the follow request
		follower, err := queries.GetFollowerByID(db, followerRowID)
		if err != nil {
			http.Error(w, "follow request not found", http.StatusNotFound)
			return
		}
		if follower.FollowingID != currentUserID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		if err := queries.AcceptFollowRequest(db, followerRowID); err != nil {
			http.Error(w, "failed to accept follow request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
	}
}

// DeclineFollowHandler handles PUT /api/followers/:id/decline
func DeclineFollowHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// extract :id from /api/followers/:id/decline
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "invalid URL", http.StatusBadRequest)
			return
		}
		followerRowID, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "invalid follower ID", http.StatusBadRequest)
			return
		}

		// verify the current user is the recipient of the follow request
		follower, err := queries.GetFollowerByID(db, followerRowID)
		if err != nil {
			http.Error(w, "follow request not found", http.StatusNotFound)
			return
		}
		if follower.FollowingID != currentUserID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		if err := queries.DeclineFollowRequest(db, followerRowID); err != nil {
			http.Error(w, "failed to decline follow request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "declined"})
	}
}

// UnfollowHandler handles POST /api/users/:id/unfollow
func UnfollowHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// extract :id from /api/users/:id/unfollow
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

		if err := queries.Unfollow(db, currentUserID, targetUserID); err != nil {
			http.Error(w, "failed to unfollow user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "unfollowed"})
	}
}

// GetFollowersHandler handles GET /api/users/:id/followers
func GetFollowersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// extract :id from /api/users/:id/followers
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "invalid URL", http.StatusBadRequest)
			return
		}
		userID, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		followers, err := queries.GetFollowers(db, userID)
		if err != nil {
			http.Error(w, "failed to fetch followers", http.StatusInternalServerError)
			return
		}

		if followers == nil {
			followers = []models.Follower{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(followers)
	}
}

// GetFollowingHandler handles GET /api/users/:id/following
func GetFollowingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// extract :id from /api/users/:id/following
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "invalid URL", http.StatusBadRequest)
			return
		}
		userID, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		following, err := queries.GetFollowing(db, userID)
		if err != nil {
			http.Error(w, "failed to fetch following", http.StatusInternalServerError)
			return
		}

		if following == nil {
			following = []models.Follower{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(following)
	}
}
