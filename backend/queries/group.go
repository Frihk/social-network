package queries

import (
	"social/pkg/db/sqlite"
)

// IsGroupMember returns true if the user is an accepted member of the group
func IsGroupMember(userID, groupID string) bool {
	query := `
		SELECT COUNT(*) > 0
		FROM group_members
		WHERE user_id = ? AND group_id = ? AND status = 'accepted'
	`
	var exists bool
	err := sqlite.DB.QueryRow(query, userID, groupID).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// GetUserGroups returns all groups the user is a member of
func GetUserGroups(userID string) ([]map[string]interface{}, error) {
	query := `
		SELECT g.id, g.name, g.description, COUNT(gm.user_id) as member_count
		FROM groups g
		JOIN group_members gm ON gm.group_id = g.id
		WHERE gm.user_id = ? AND gm.status = 'accepted'
		GROUP BY g.id, g.name, g.description
		ORDER BY g.name ASC
	`
	rows, err := sqlite.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []map[string]interface{}
	for rows.Next() {
		var id, name, description string
		var memberCount int
		if err := rows.Scan(&id, &name, &description, &memberCount); err != nil {
			return nil, err
		}
		group := map[string]interface{}{
			"id":           id,
			"name":         name,
			"description":  description,
			"member_count": memberCount,
		}
		groups = append(groups, group)
	}
	return groups, rows.Err()
}
