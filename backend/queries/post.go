package queries

import (
	"database/sql"
	"fmt"

	"social/models"
)

// CreatePost inserts a new post and its allowed viewers in a transaction
func CreatePost(db *sql.DB, post models.Post, allowedViewers []int) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO posts (user_id, group_id, content, image_path, privacy, created_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		post.UserID, post.GroupID, post.Content, post.ImagePath, post.Privacy,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert post: %w", err)
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get post ID: %w", err)
	}

	if post.Privacy == "private" && len(allowedViewers) > 0 {
		for _, viewerID := range allowedViewers {
			_, err := tx.Exec(`
				INSERT INTO post_allowed_viewers (post_id, user_id)
				VALUES (?, ?)`,
				postID, viewerID,
			)
			if err != nil {
				return 0, fmt.Errorf("failed to insert allowed viewer: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return postID, nil
}

// GetFeed returns posts visible to the current user ordered by newest first
func GetFeed(db *sql.DB, userID int) ([]models.Post, error) {
	rows, err := db.Query(`
		SELECT DISTINCT p.id, p.user_id, p.group_id, p.content, p.image_path, p.privacy, p.created_at
		FROM posts p
		LEFT JOIN followers f
			ON f.following_id = p.user_id
			AND f.follower_id = ?
			AND f.status = 'accepted'
		LEFT JOIN post_allowed_viewers pav
			ON pav.post_id = p.id
			AND pav.user_id = ?
		WHERE p.group_id IS NULL
		AND (
			p.user_id = ?
			OR p.privacy = 'public'
			OR (p.privacy = 'almost_private' AND f.follower_id IS NOT NULL)
			OR (p.privacy = 'private' AND pav.user_id IS NOT NULL)
		)
		ORDER BY p.created_at DESC`,
		userID, userID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}
	defer rows.Close()

	return scanPosts(rows)
}

// GetPostsByUserID returns all posts by a specific user visible to the current user
func GetPostsByUserID(db *sql.DB, targetUserID, currentUserID int) ([]models.Post, error) {
	rows, err := db.Query(`
		SELECT DISTINCT p.id, p.user_id, p.group_id, p.content, p.image_path, p.privacy, p.created_at
		FROM posts p
		LEFT JOIN followers f
			ON f.following_id = p.user_id
			AND f.follower_id = ?
			AND f.status = 'accepted'
		LEFT JOIN post_allowed_viewers pav
			ON pav.post_id = p.id
			AND pav.user_id = ?
		WHERE p.user_id = ?
		AND p.group_id IS NULL
		AND (
			p.user_id = ?
			OR p.privacy = 'public'
			OR (p.privacy = 'almost_private' AND f.follower_id IS NOT NULL)
			OR (p.privacy = 'private' AND pav.user_id IS NOT NULL)
		)
		ORDER BY p.created_at DESC`,
		currentUserID, currentUserID, targetUserID, currentUserID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user posts: %w", err)
	}
	defer rows.Close()

	return scanPosts(rows)
}

// GetPostByID returns a single post by its ID
func GetPostByID(db *sql.DB, postID int) (models.Post, error) {
	var post models.Post
	err := db.QueryRow(`
		SELECT id, user_id, group_id, content, image_path, privacy, created_at
		FROM posts
		WHERE id = ?`, postID,
	).Scan(
		&post.ID, &post.UserID, &post.GroupID,
		&post.Content, &post.ImagePath,
		&post.Privacy, &post.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return post, fmt.Errorf("post not found")
	}
	if err != nil {
		return post, fmt.Errorf("failed to fetch post: %w", err)
	}
	return post, nil
}

// CreateComment inserts a new comment on a post
func CreateComment(db *sql.DB, comment models.Comment) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO comments (post_id, user_id, content, image_path, created_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		comment.PostID, comment.UserID, comment.Content, comment.ImagePath,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert comment: %w", err)
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get comment ID: %w", err)
	}

	return commentID, nil
}

// GetCommentsByPostID returns all comments for a given post
func GetCommentsByPostID(db *sql.DB, postID int) ([]models.Comment, error) {
	rows, err := db.Query(`
		SELECT id, post_id, user_id, content, image_path, created_at
		FROM comments
		WHERE post_id = ?
		ORDER BY created_at ASC`, postID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(
			&c.ID, &c.PostID, &c.UserID,
			&c.Content, &c.ImagePath, &c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, c)
	}

	return comments, nil
}

// scanPosts is a helper that scans rows into a slice of Post
func scanPosts(rows *sql.Rows) ([]models.Post, error) {
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.GroupID,
			&p.Content, &p.ImagePath,
			&p.Privacy, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, nil
}
