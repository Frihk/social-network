package models

import "time"

type Group struct {
	ID          string    `json:"id"`
	CreatorID   string    `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupMember struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Event struct {
	ID          string    `json:"id"`
	GroupID     string    `json:"group_id"`
	CreatorID   string    `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	CreatedAt   time.Time `json:"created_at"`
}

type EventResponse struct {
	ID        string    `json:"id"`
	EventID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	Response  string    `json:"response"`
	CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	RelatedID string    `json:"related_id"`
	Message   string    `json:"message"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateGroupRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type InviteUserRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type RequestToJoinRequest struct {
}

type RespondToInviteRequest struct {
}

type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time" binding:"required"`
}

type RespondToEventRequest struct {
	Response string `json:"response" binding:"required,oneof=going not_going"`
}

type GroupDetailResponse struct {
	Group         *Group         `json:"group"`
	Members       []*GroupMember `json:"members"`
	MemberCount   int            `json:"member_count"`
	AcceptedCount int            `json:"accepted_count"`
}

type EventDetailResponse struct {
	Event     *Event           `json:"event"`
	Responses []EventResponse  `json:"responses"`
	GoingCount    int          `json:"going_count"`
	NotGoingCount int          `json:"not_going_count"`
}

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

type Reaction struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
}

type ReactionRequest struct {
	Emoji string `json:"emoji"`
}