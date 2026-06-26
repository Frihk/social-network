package queries

import (
	"database/sql"
	"social/pkg/db/sqlite"
)

// GetDMEligibleUsers returns all users the logged-in user can DM
func GetDMEligibleUsers(userID string) ([]map[string]interface{}, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, u.avatar
		FROM users u
		JOIN followers f ON (u.id = f.follower_id OR u.id = f.following_id)
		WHERE (f.follower_id = ? OR f.following_id = ?) 
		AND f.status = 'accepted'
		AND u.id != ?
	`
	rows, err := sqlite.DB.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	seen := make(map[string]bool)
	for rows.Next() {
		var id, firstName, lastName string
		var avatar sql.NullString
		if err := rows.Scan(&id, &firstName, &lastName, &avatar); err != nil {
			return nil, err
		}
		if !seen[id] {
			seen[id] = true
			u := map[string]interface{}{
				"id":         id,
				"first_name": firstName,
				"last_name":  lastName,
				"avatar":     nil,
			}
			if avatar.Valid {
				u["avatar"] = avatar.String
			}
			users = append(users, u)
		}
	}
	// Return empty slice instead of nil for JSON serialization
	if users == nil {
		users = []map[string]interface{}{}
	}
	return users, nil
}

func CanDM(userID1, userID2 string) bool {
	query := `
		SELECT COUNT(*) 
		FROM followers 
		WHERE status = 'accepted' AND 
		((follower_id = ? AND following_id = ?) OR (follower_id = ? AND following_id = ?))
	`
	var count int
	err := sqlite.DB.QueryRow(query, userID1, userID2, userID2, userID1).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func IsGroupMember(userID, groupID string) bool {
	query := `
		SELECT COUNT(*) 
		FROM group_members 
		WHERE group_id = ? AND user_id = ? AND status = 'accepted'
	`
	var count int
	err := sqlite.DB.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func GetUserGroups(userID string) ([]map[string]interface{}, error) {
	query := `
		SELECT g.id, g.title, g.description, gm.status
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = ?
	`
	rows, err := sqlite.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []map[string]interface{}
	for rows.Next() {
		var id, title, description, status string
		if err := rows.Scan(&id, &title, &description, &status); err != nil {
			return nil, err
		}
		g := map[string]interface{}{
			"id":          id,
			"title":       title,
			"description": description,
			"status":      status,
		}
		groups = append(groups, g)
	}
	if groups == nil {
		groups = []map[string]interface{}{}
	}
	return groups, nil
}
