package queries

import (
	"database/sql"
	"social/pkg/db/sqlite"
	"github.com/google/uuid"
)

func SavePrivateMessage(senderID, receiverID, content string) error {
	query := `INSERT INTO messages (id, sender_id, receiver_id, content) VALUES (?, ?, ?, ?)`
	_, err := sqlite.DB.Exec(query, uuid.New().String(), senderID, receiverID, content)
	return err
}

func SaveGroupMessage(groupID, senderID, content string) error {
	query := `INSERT INTO group_messages (id, group_id, sender_id, content) VALUES (?, ?, ?, ?)`
	_, err := sqlite.DB.Exec(query, uuid.New().String(), groupID, senderID, content)
	return err
}

func GetPrivateMessageHistory(userID1, userID2 string) ([]map[string]interface{}, error) {
	query := `
		SELECT id, sender_id, receiver_id, content, created_at
		FROM messages
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
		ORDER BY created_at ASC
	`
	rows, err := sqlite.DB.Query(query, userID1, userID2, userID2, userID1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var id, senderID, receiverID, content string
		var createdAt sql.NullString
		if err := rows.Scan(&id, &senderID, &receiverID, &content, &createdAt); err != nil {
			return nil, err
		}
		msg := map[string]interface{}{
			"id":          id,
			"sender_id":   senderID,
			"receiver_id": receiverID,
			"content":     content,
			"created_at":  createdAt.String,
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

func GetGroupMessageHistory(groupID string) ([]map[string]interface{}, error) {
	query := `
		SELECT id, group_id, sender_id, content, created_at
		FROM group_messages
		WHERE group_id = ?
		ORDER BY created_at ASC
	`
	rows, err := sqlite.DB.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var id, groupIDCol, senderID, content string
		var createdAt sql.NullString
		if err := rows.Scan(&id, &groupIDCol, &senderID, &content, &createdAt); err != nil {
			return nil, err
		}
		msg := map[string]interface{}{
			"id":         id,
			"group_id":   groupIDCol,
			"sender_id":  senderID,
			"content":    content,
			"created_at": createdAt.String,
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}
