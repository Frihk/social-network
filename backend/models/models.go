package models

import (
	"time"
)


type Post struct {
	ID 	  int	    `json:"id"`
	UserID	  int	    `json:"user_id"`
	GroupID	  *int	    `json:"group_id"`
	Content	  string    `json:"content"`
	ImagePath *string   `json:"image_path"`
	Privacy   string    `json:"privacy"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID	  int       `json:"id"`
	PostID	  int 	    `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	ImagePath *string   `json:"image_path"`
	CreatedAt time.Time `json:"created_at"`
}

type Follower struct {
	ID          int       `json:"id"`
	FollowerID  int       `json:"follower_id"`
	FollowingID int       `json:"following_id"`
	Status 	    string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

