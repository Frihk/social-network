package queries

import (
	"database/sql"
	"fmt"

	"social/models"
)

// CreateFollower inserts a new follower relationship with the given status
func CreateFollower(db *sql.DB, followerID, followingID int, status string) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO followers (follower_id, following_id, status, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)`,
		followerID, followingID, status,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create follower: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get follower ID: %w", err)
	}

	return id, nil
}

// GetFollowerByID returns a single follower row by its ID
func GetFollowerByID(db *sql.DB, id int) (models.Follower, error) {
	var f models.Follower
	err := db.QueryRow(`
		SELECT id, follower_id, following_id, status, created_at
		FROM followers
		WHERE id = ?`, id,
	).Scan(&f.ID, &f.FollowerID, &f.FollowingID, &f.Status, &f.CreatedAt)
	if err == sql.ErrNoRows {
		return f, fmt.Errorf("follow request not found")
	}
	if err != nil {
		return f, fmt.Errorf("failed to fetch follower: %w", err)
	}
	return f, nil
}

// GetFollowStatus returns the follow relationship between two users if it exists
func GetFollowStatus(db *sql.DB, followerID, followingID int) (models.Follower, error) {
	var f models.Follower
	err := db.QueryRow(`
		SELECT id, follower_id, following_id, status, created_at
		FROM followers
		WHERE follower_id = ? AND following_id = ?`,
		followerID, followingID,
	).Scan(&f.ID, &f.FollowerID, &f.FollowingID, &f.Status, &f.CreatedAt)
	if err == sql.ErrNoRows {
		return f, fmt.Errorf("no follow relationship found")
	}
	if err != nil {
		return f, fmt.Errorf("failed to fetch follow status: %w", err)
	}
	return f, nil
}

// AcceptFollowRequest updates a follow request status to accepted
func AcceptFollowRequest(db *sql.DB, id int) error {
	result, err := db.Exec(`
		UPDATE followers
		SET status = 'accepted'
		WHERE id = ?`, id,
	)
	if err != nil {
		return fmt.Errorf("failed to accept follow request: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("follow request not found")
	}

	return nil
}

// DeclineFollowRequest deletes a follow request by its ID
func DeclineFollowRequest(db *sql.DB, id int) error {
	result, err := db.Exec(`
		DELETE FROM followers
		WHERE id = ?`, id,
	)
	if err != nil {
		return fmt.Errorf("failed to decline follow request: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("follow request not found")
	}

	return nil
}

// Unfollow deletes a follow relationship by follower and following IDs
func Unfollow(db *sql.DB, followerID, followingID int) error {
	result, err := db.Exec(`
		DELETE FROM followers
		WHERE follower_id = ? AND following_id = ?`,
		followerID, followingID,
	)
	if err != nil {
		return fmt.Errorf("failed to unfollow: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("follow relationship not found")
	}

	return nil
}

// GetFollowers returns all accepted followers for a given user
func GetFollowers(db *sql.DB, userID int) ([]models.Follower, error) {
	rows, err := db.Query(`
		SELECT id, follower_id, following_id, status, created_at
		FROM followers
		WHERE following_id = ? AND status = 'accepted'
		ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch followers: %w", err)
	}
	defer rows.Close()

	return scanFollowers(rows)
}

// GetFollowing returns all users a given user is following
func GetFollowing(db *sql.DB, userID int) ([]models.Follower, error) {
	rows, err := db.Query(`
		SELECT id, follower_id, following_id, status, created_at
		FROM followers
		WHERE follower_id = ? AND status = 'accepted'
		ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch following: %w", err)
	}
	defer rows.Close()

	return scanFollowers(rows)
}

// scanFollowers is a helper that scans rows into a slice of Follower
func scanFollowers(rows *sql.Rows) ([]models.Follower, error) {
	var followers []models.Follower
	for rows.Next() {
		var f models.Follower
		if err := rows.Scan(
			&f.ID, &f.FollowerID, &f.FollowingID,
			&f.Status, &f.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan follower: %w", err)
		}
		followers = append(followers, f)
	}
	return followers, nil
}
