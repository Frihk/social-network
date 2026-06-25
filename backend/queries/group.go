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
