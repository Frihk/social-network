package queries

import (
	"database/sql"

	"github.com/google/uuid"
	"social/pkg/db/sqlite"
)

func SavePrivateMessage(senderID, receiverID, content string) (string, error) {
	id := uuid.New().String()
	query := `INSERT INTO messages (id, sender_id, receiver_id, content) VALUES (?, ?, ?, ?)`
	_, err := sqlite.DB.Exec(query, id, senderID, receiverID, content)
	return id, err
}

func SaveGroupMessage(groupID, senderID, content string) (string, error) {
	id := uuid.New().String()
	query := `INSERT INTO group_messages (id, group_id, sender_id, content) VALUES (?, ?, ?, ?)`
	_, err := sqlite.DB.Exec(query, id, groupID, senderID, content)
	return id, err
}

func GetAcceptedGroupMemberIDs(groupID string) ([]string, error) {
	rows, err := sqlite.DB.Query(`
		SELECT user_id
		FROM group_members
		WHERE group_id = ? AND status = 'accepted'
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		memberIDs = append(memberIDs, userID)
	}
	return memberIDs, rows.Err()
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
