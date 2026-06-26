package middleware

import (
	"context"
	"net/http"
	"time"

	"social/models"
	"social/pkg/db/sqlite"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		session, err := models.GetSessionByID(sqlite.DB, cookie.Value)
		if err != nil {
			http.Error(w, `{"error":"Invalid session"}`, http.StatusUnauthorized)
			return
		}

		if time.Now().After(session.ExpiresAt) {
			models.DeleteSession(sqlite.DB, session.ID)
			http.Error(w, `{"error":"Session expired"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}
