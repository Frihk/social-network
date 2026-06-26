package queries

import (
	"database/sql"
	"social/models"
)

func GetFeed(db *sql.DB, userID string) ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.content, p.privacy, p.image_path, p.created_at,
		       u.first_name || ' ' || u.last_name AS author_name,
		       u.avatar AS author_avatar
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.privacy = 'public'
		   OR p.user_id = ?
		   OR (p.privacy = 'almost_private' AND EXISTS (
		       SELECT 1 FROM followers f
		       WHERE f.follower_id = ? AND f.following_id = p.user_id AND f.status = 'accepted'
		   ))
		   OR (p.privacy = 'private' AND EXISTS (
		       SELECT 1 FROM post_allowed_viewers pav
		       WHERE pav.post_id = p.id AND pav.user_id = ?
		   ))
		ORDER BY p.created_at DESC
	`
	rows, err := db.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.Privacy, &p.ImagePath, &p.CreatedAt, &p.AuthorName, &p.AuthorAvatar); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func CreatePost(db *sql.DB, post models.Post, allowedViewers []string) (int64, error) {
	query := `INSERT INTO posts (user_id, content, privacy, image_path) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, post.UserID, post.Content, post.Privacy, post.ImagePath)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if post.Privacy == "private" && len(allowedViewers) > 0 {
		for _, viewerID := range allowedViewers {
			_, err := db.Exec(`INSERT INTO post_allowed_viewers (post_id, user_id) VALUES (?, ?)`, postID, viewerID)
			if err != nil {
				return 0, err
			}
		}
	}

	return postID, nil
}

func GetPostsByUserID(db *sql.DB, targetUserID string, viewerID string) ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, p.content, p.privacy, p.image_path, p.created_at,
		       u.first_name || ' ' || u.last_name AS author_name,
		       u.avatar AS author_avatar
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.user_id = ?
		  AND (
		       p.privacy = 'public'
		    OR p.user_id = ?
		    OR (p.privacy = 'almost_private' AND EXISTS (
		        SELECT 1 FROM followers f
		        WHERE f.follower_id = ? AND f.following_id = p.user_id AND f.status = 'accepted'
		    ))
		    OR (p.privacy = 'private' AND EXISTS (
		        SELECT 1 FROM post_allowed_viewers pav
		        WHERE pav.post_id = p.id AND pav.user_id = ?
		    ))
		  )
		ORDER BY p.created_at DESC
	`
	rows, err := db.Query(query, targetUserID, viewerID, viewerID, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.Privacy, &p.ImagePath, &p.CreatedAt, &p.AuthorName, &p.AuthorAvatar); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func CreateComment(db *sql.DB, comment models.Comment) (int64, error) {
	query := `INSERT INTO comments (post_id, user_id, content, image_path) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, comment.PostID, comment.UserID, comment.Content, comment.ImagePath)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetCommentsByPostID(db *sql.DB, postID int) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.image_path, c.created_at,
		       u.first_name || ' ' || u.last_name AS author_name,
		       u.avatar AS author_avatar
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.ImagePath, &c.CreatedAt, &c.AuthorName, &c.AuthorAvatar); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
