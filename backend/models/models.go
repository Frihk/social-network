package models

import "time"

// Group represents a social group
type Group struct {
	ID          string    `json:"id"`
	CreatorID   string    `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// GroupMember represents a user's membership in a group
type GroupMember struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"` // 'invited', 'requested', 'accepted'
	CreatedAt time.Time `json:"created_at"`
}

// Event represents a group event
type Event struct {
	ID          string    `json:"id"`
	GroupID     string    `json:"group_id"`
	CreatorID   string    `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	CreatedAt   time.Time `json:"created_at"`
}

// EventResponse represents a user's RSVP to an event
type EventResponse struct {
	ID        string    `json:"id"`
	EventID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	Response  string    `json:"response"` // 'going', 'not_going'
	CreatedAt time.Time `json:"created_at"`
}

// Notification represents a notification to a user
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"` // 'group_invite', 'group_request', 'event_created', etc.
	RelatedID string    `json:"related_id"` // ID of related group, event, or user
	Message   string    `json:"message"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateGroupRequest is the request body for creating a group
type CreateGroupRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// InviteUserRequest is the request body for inviting a user to a group
type InviteUserRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// RequestToJoinRequest is the request body for requesting to join a group
type RequestToJoinRequest struct {
	// No additional fields needed beyond the authenticated user
}

// RespondToInviteRequest is used for accepting/declining invites
type RespondToInviteRequest struct {
	// No additional fields needed
}

// CreateEventRequest is the request body for creating an event
type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time" binding:"required"`
}

// RespondToEventRequest is the request body for responding to an event
type RespondToEventRequest struct {
	Response string `json:"response" binding:"required,oneof=going not_going"`
}

// GroupDetailResponse is the response for GET /api/groups/:id (members only)
type GroupDetailResponse struct {
	Group         *Group         `json:"group"`
	Members       []*GroupMember `json:"members"`
	MemberCount   int            `json:"member_count"`
	AcceptedCount int            `json:"accepted_count"`
}

// EventDetailResponse includes event and responses
type EventDetailResponse struct {
	Event     *Event           `json:"event"`
	Responses []EventResponse  `json:"responses"`
	GoingCount    int          `json:"going_count"`
	NotGoingCount int          `json:"not_going_count"`
}

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
