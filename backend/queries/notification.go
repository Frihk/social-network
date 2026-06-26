package queries

import (
	"database/sql"
	"social/pkg/db/sqlite"
	"github.com/google/uuid"
)

func CreateNotification(userID, notificationType, actorID, entityID string) error {
	query := `INSERT INTO notifications (id, user_id, type, actor_id, entity_id, is_read) VALUES (?, ?, ?, ?, ?, 0)`
	_, err := sqlite.DB.Exec(query, uuid.New().String(), userID, notificationType, actorID, entityID)
	return err
}

func GetNotificationsForUser(userID string) ([]map[string]interface{}, error) {
	query := `
		SELECT id, user_id, type, actor_id, entity_id, is_read, created_at
		FROM notifications
		WHERE user_id = ?
		ORDER BY is_read ASC, created_at DESC
	`
	rows, err := sqlite.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []map[string]interface{}
	for rows.Next() {
		var id, userID, notifType, actorID, entityID string
		var isRead int
		var createdAt sql.NullString
		if err := rows.Scan(&id, &userID, &notifType, &actorID, &entityID, &isRead, &createdAt); err != nil {
			return nil, err
		}
		notif := map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"type":       notifType,
			"actor_id":   actorID,
			"entity_id":  entityID,
			"is_read":    isRead,
			"created_at": createdAt.String,
		}
		notifications = append(notifications, notif)
	}
	return notifications, rows.Err()
}

func MarkAsRead(notificationID string) error {
	query := `UPDATE notifications SET is_read = 1 WHERE id = ?`
	_, err := sqlite.DB.Exec(query, notificationID)
	return err
}

func MarkAllAsRead(userID string) error {
	query := `UPDATE notifications SET is_read = 1 WHERE user_id = ?`
	_, err := sqlite.DB.Exec(query, userID)
	return err
}

// GetNotificationOwner returns the user_id of a notification, or empty string if not found
func GetNotificationOwner(notificationID string) (string, error) {
	var userID string
	err := sqlite.DB.QueryRow(`SELECT user_id FROM notifications WHERE id = ?`, notificationID).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}
