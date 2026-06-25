package queries

import (
	"social/pkg/db/sqlite"
)

// CanDM returns true if at least one user follows the other with status = 'accepted'
func CanDM(userID1, userID2 string) bool {
	query := `
		SELECT COUNT(*) > 0
		FROM followers
		WHERE (follower_id = ? AND following_id = ? AND status = 'accepted')
		   OR (follower_id = ? AND following_id = ? AND status = 'accepted')
	`
	var exists bool
	err := sqlite.DB.QueryRow(query, userID1, userID2, userID2, userID1).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// GetDMEligibleUsers returns all users where the logged-in user can DM
// (either follows them or they follow the logged-in user with accepted status)
func GetDMEligibleUsers(loggedInUserID string) ([]map[string]interface{}, error) {
	query := `
		SELECT DISTINCT u.id, u.first_name, u.last_name, u.avatar, u.nickname
		FROM users u
		JOIN followers f1 ON f1.following_id = u.id AND f1.follower_id = ? AND f1.status = 'accepted'
		UNION
		SELECT DISTINCT u.id, u.first_name, u.last_name, u.avatar, u.nickname
		FROM users u
		JOIN followers f2 ON f2.follower_id = u.id AND f2.following_id = ? AND f2.status = 'accepted'
	`
	rows, err := sqlite.DB.Query(query, loggedInUserID, loggedInUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id, firstName, lastName string
		var avatar, nickname *string
		if err := rows.Scan(&id, &firstName, &lastName, &avatar, &nickname); err != nil {
			return nil, err
		}
		user := map[string]interface{}{
			"id":          id,
			"first_name":  firstName,
			"last_name":   lastName,
			"avatar_path": avatar,
			"nickname":    nickname,
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
