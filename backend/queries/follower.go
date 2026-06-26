package queries

import (
	"database/sql"
	"social/models"
)

func GetFollowStatus(db *sql.DB, followerID string, followingID string) (*models.Follower, error) {
	f := &models.Follower{}
	query := `SELECT id, follower_id, following_id, status, created_at FROM followers WHERE follower_id = ? AND following_id = ?`
	err := db.QueryRow(query, followerID, followingID).Scan(&f.ID, &f.FollowerID, &f.FollowingID, &f.Status, &f.CreatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func CreateFollower(db *sql.DB, followerID string, followingID string, status string) (int64, error) {
	query := `INSERT INTO followers (follower_id, following_id, status) VALUES (?, ?, ?)`
	result, err := db.Exec(query, followerID, followingID, status)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetFollowerByID(db *sql.DB, id int) (*models.Follower, error) {
	f := &models.Follower{}
	query := `SELECT id, follower_id, following_id, status, created_at FROM followers WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&f.ID, &f.FollowerID, &f.FollowingID, &f.Status, &f.CreatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func AcceptFollowRequest(db *sql.DB, id int) error {
	_, err := db.Exec(`UPDATE followers SET status = 'accepted' WHERE id = ?`, id)
	return err
}

func DeclineFollowRequest(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM followers WHERE id = ?`, id)
	return err
}

func Unfollow(db *sql.DB, followerID string, followingID string) error {
	_, err := db.Exec(`DELETE FROM followers WHERE follower_id = ? AND following_id = ?`, followerID, followingID)
	return err
}

func GetFollowers(db *sql.DB, userID string) ([]models.Follower, error) {
	rows, err := db.Query(
		`SELECT id, follower_id, following_id, status, created_at FROM followers WHERE following_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []models.Follower
	for rows.Next() {
		var f models.Follower
		if err := rows.Scan(&f.ID, &f.FollowerID, &f.FollowingID, &f.Status, &f.CreatedAt); err != nil {
			return nil, err
		}
		followers = append(followers, f)
	}
	return followers, nil
}

func GetFollowing(db *sql.DB, userID string) ([]models.Follower, error) {
	rows, err := db.Query(
		`SELECT id, follower_id, following_id, status, created_at FROM followers WHERE follower_id = ? AND status = 'accepted'`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []models.Follower
	for rows.Next() {
		var f models.Follower
		if err := rows.Scan(&f.ID, &f.FollowerID, &f.FollowingID, &f.Status, &f.CreatedAt); err != nil {
			return nil, err
		}
		following = append(following, f)
	}
	return following, nil
}
