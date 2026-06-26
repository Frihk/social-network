package models

import "time"

type Post struct {
	ID           int64     `json:"id"`
	UserID       string    `json:"user_id"`
	Content      string    `json:"content"`
	Privacy      string    `json:"privacy"`
	ImagePath    *string   `json:"image_path"`
	CreatedAt    time.Time `json:"created_at"`
	AuthorName   string    `json:"author_name,omitempty"`
	AuthorAvatar *string   `json:"author_avatar,omitempty"`
}

type Comment struct {
	ID           int64     `json:"id"`
	PostID       int       `json:"post_id"`
	UserID       string    `json:"user_id"`
	Content      string    `json:"content"`
	ImagePath    *string   `json:"image_path"`
	CreatedAt    time.Time `json:"created_at"`
	AuthorName   string    `json:"author_name,omitempty"`
	AuthorAvatar *string   `json:"author_avatar,omitempty"`
}

type Follower struct {
	ID          int       `json:"id"`
	FollowerID  string    `json:"follower_id"`
	FollowingID string    `json:"following_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
