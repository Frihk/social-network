package queries

import (
	"database/sql"
	"errors"
	"social/models"
	"time"

	"github.com/google/uuid"
)

// GroupQueries contains all group-related database operations
type GroupQueries struct {
	db *sql.DB
}

// NewGroupQueries creates a new GroupQueries instance
func NewGroupQueries(db *sql.DB) *GroupQueries {
	return &GroupQueries{db: db}
}

// CreateGroup creates a new group and sets the creator as the first accepted member
func (q *GroupQueries) CreateGroup(title, description, creatorID string) (*models.Group, error) {
	groupID := uuid.New().String()

	tx, err := q.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create group
	group := &models.Group{
		ID:          groupID,
		CreatorID:   creatorID,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
	}

	_, err = tx.Exec(
		"INSERT INTO groups (id, creator_id, title, description, created_at) VALUES (?, ?, ?, ?, ?)",
		group.ID, group.CreatorID, group.Title, group.Description, group.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Add creator as accepted member
	memberID := uuid.New().String()
	_, err = tx.Exec(
		"INSERT INTO group_members (id, group_id, user_id, status, created_at) VALUES (?, ?, ?, ?, ?)",
		memberID, groupID, creatorID, "accepted", time.Now(),
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return group, nil
}

// GetAllGroups returns all groups
func (q *GroupQueries) GetAllGroups() ([]*models.Group, error) {
	rows, err := q.db.Query("SELECT id, creator_id, title, description, created_at FROM groups ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		group := &models.Group{}
		if err := rows.Scan(&group.ID, &group.CreatorID, &group.Title, &group.Description, &group.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// GetGroupByID returns a group by ID
func (q *GroupQueries) GetGroupByID(groupID string) (*models.Group, error) {
	group := &models.Group{}
	err := q.db.QueryRow(
		"SELECT id, creator_id, title, description, created_at FROM groups WHERE id = ?",
		groupID,
	).Scan(&group.ID, &group.CreatorID, &group.Title, &group.Description, &group.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("group not found")
		}
		return nil, err
	}

	return group, nil
}

// GetGroupMembers returns all members of a group
func (q *GroupQueries) GetGroupMembers(groupID string) ([]*models.GroupMember, error) {
	rows, err := q.db.Query(
		"SELECT id, group_id, user_id, status, created_at FROM group_members WHERE group_id = ? ORDER BY created_at ASC",
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.GroupMember
	for rows.Next() {
		member := &models.GroupMember{}
		if err := rows.Scan(&member.ID, &member.GroupID, &member.UserID, &member.Status, &member.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

// GetAcceptedMembers returns only accepted members of a group
func (q *GroupQueries) GetAcceptedMembers(groupID string) ([]*models.GroupMember, error) {
	rows, err := q.db.Query(
		"SELECT id, group_id, user_id, status, created_at FROM group_members WHERE group_id = ? AND status = 'accepted' ORDER BY created_at ASC",
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.GroupMember
	for rows.Next() {
		member := &models.GroupMember{}
		if err := rows.Scan(&member.ID, &member.GroupID, &member.UserID, &member.Status, &member.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

// GetMembershipStatus returns the membership status of a user in a group
func (q *GroupQueries) GetMembershipStatus(groupID, userID string) (string, error) {
	var status string
	err := q.db.QueryRow(
		"SELECT status FROM group_members WHERE group_id = ? AND user_id = ?",
		groupID, userID,
	).Scan(&status)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Not a member
		}
		return "", err
	}

	return status, nil
}

// InviteUserToGroup creates an invited membership entry
func (q *GroupQueries) InviteUserToGroup(groupID, userID string) (*models.GroupMember, error) {
	memberID := uuid.New().String()
	now := time.Now()

	member := &models.GroupMember{
		ID:        memberID,
		GroupID:   groupID,
		UserID:    userID,
		Status:    "invited",
		CreatedAt: now,
	}

	_, err := q.db.Exec(
		"INSERT INTO group_members (id, group_id, user_id, status, created_at) VALUES (?, ?, ?, ?, ?)",
		member.ID, member.GroupID, member.UserID, member.Status, member.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return member, nil
}

// RequestToJoinGroup creates a requested membership entry
func (q *GroupQueries) RequestToJoinGroup(groupID, userID string) (*models.GroupMember, error) {
	memberID := uuid.New().String()
	now := time.Now()

	member := &models.GroupMember{
		ID:        memberID,
		GroupID:   groupID,
		UserID:    userID,
		Status:    "requested",
		CreatedAt: now,
	}

	_, err := q.db.Exec(
		"INSERT INTO group_members (id, group_id, user_id, status, created_at) VALUES (?, ?, ?, ?, ?)",
		member.ID, member.GroupID, member.UserID, member.Status, member.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return member, nil
}

// AcceptMembership sets a membership status to accepted
func (q *GroupQueries) AcceptMembership(groupID, userID string) error {
	result, err := q.db.Exec(
		"UPDATE group_members SET status = 'accepted' WHERE group_id = ? AND user_id = ?",
		groupID, userID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("membership not found")
	}

	return nil
}

// DeclineMembership deletes a membership entry
func (q *GroupQueries) DeclineMembership(groupID, userID string) error {
	result, err := q.db.Exec(
		"DELETE FROM group_members WHERE group_id = ? AND user_id = ?",
		groupID, userID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("membership not found")
	}

	return nil
}

// CreateEvent creates a new event in a group
func (q *GroupQueries) CreateEvent(groupID, creatorID, title, description string, eventTime time.Time) (*models.Event, error) {
	eventID := uuid.New().String()

	event := &models.Event{
		ID:          eventID,
		GroupID:     groupID,
		CreatorID:   creatorID,
		Title:       title,
		Description: description,
		EventTime:   eventTime,
		CreatedAt:   time.Now(),
	}

	_, err := q.db.Exec(
		"INSERT INTO events (id, group_id, creator_id, title, description, event_time, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		event.ID, event.GroupID, event.CreatorID, event.Title, event.Description, event.EventTime, event.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// GetEventByID returns an event by ID
func (q *GroupQueries) GetEventByID(eventID string) (*models.Event, error) {
	event := &models.Event{}
	err := q.db.QueryRow(
		"SELECT id, group_id, creator_id, title, description, event_time, created_at FROM events WHERE id = ?",
		eventID,
	).Scan(&event.ID, &event.GroupID, &event.CreatorID, &event.Title, &event.Description, &event.EventTime, &event.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("event not found")
		}
		return nil, err
	}

	return event, nil
}

// GetGroupEvents returns all events for a group
func (q *GroupQueries) GetGroupEvents(groupID string) ([]*models.Event, error) {
	rows, err := q.db.Query(
		"SELECT id, group_id, creator_id, title, description, event_time, created_at FROM events WHERE group_id = ? ORDER BY event_time DESC",
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		event := &models.Event{}
		if err := rows.Scan(&event.ID, &event.GroupID, &event.CreatorID, &event.Title, &event.Description, &event.EventTime, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// RespondToEvent creates or updates a response to an event
func (q *GroupQueries) RespondToEvent(eventID, userID, response string) (*models.EventResponse, error) {
	// First check if response already exists
	var existingID string
	err := q.db.QueryRow(
		"SELECT id FROM event_responses WHERE event_id = ? AND user_id = ?",
		eventID, userID,
	).Scan(&existingID)

	if err == nil {
		// Update existing response
		_, err = q.db.Exec(
			"UPDATE event_responses SET response = ? WHERE event_id = ? AND user_id = ?",
			response, eventID, userID,
		)
		if err != nil {
			return nil, err
		}

		return &models.EventResponse{
			ID:        existingID,
			EventID:   eventID,
			UserID:    userID,
			Response:  response,
			CreatedAt: time.Now(),
		}, nil
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	// Create new response
	responseID := uuid.New().String()
	now := time.Now()

	eventResp := &models.EventResponse{
		ID:        responseID,
		EventID:   eventID,
		UserID:    userID,
		Response:  response,
		CreatedAt: now,
	}

	_, err = q.db.Exec(
		"INSERT INTO event_responses (id, event_id, user_id, response, created_at) VALUES (?, ?, ?, ?, ?)",
		eventResp.ID, eventResp.EventID, eventResp.UserID, eventResp.Response, eventResp.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return eventResp, nil
}

// GetEventResponses returns all responses for an event
func (q *GroupQueries) GetEventResponses(eventID string) ([]*models.EventResponse, error) {
	rows, err := q.db.Query(
		"SELECT id, event_id, user_id, response, created_at FROM event_responses WHERE event_id = ? ORDER BY created_at ASC",
		eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []*models.EventResponse
	for rows.Next() {
		resp := &models.EventResponse{}
		if err := rows.Scan(&resp.ID, &resp.EventID, &resp.UserID, &resp.Response, &resp.CreatedAt); err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return responses, nil
}

// GetEventResponseCounts returns the count of going/not_going responses
func (q *GroupQueries) GetEventResponseCounts(eventID string) (going, notGoing int, err error) {
	err = q.db.QueryRow(
		"SELECT COUNT(CASE WHEN response = 'going' THEN 1 END), COUNT(CASE WHEN response = 'not_going' THEN 1 END) FROM event_responses WHERE event_id = ?",
		eventID,
	).Scan(&going, &notGoing)

	return
}

// CreateNotification creates a notification for a user
func (q *GroupQueries) CreateNotification(userID, notificationType, relatedID, message string) (*models.Notification, error) {
	notifID := uuid.New().String()
	now := time.Now()

	notification := &models.Notification{
		ID:        notifID,
		UserID:    userID,
		Type:      notificationType,
		RelatedID: relatedID,
		Message:   message,
		Read:      false,
		CreatedAt: now,
	}

	_, err := q.db.Exec(
		"INSERT INTO notifications (id, user_id, type, actor_id, entity_id, is_read, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		notification.ID, notification.UserID, notification.Type, notification.UserID, notification.RelatedID, notification.Read, notification.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	dispatchNotification(userID, map[string]interface{}{
		"id":         notification.ID,
		"user_id":    notification.UserID,
		"type":       notification.Type,
		"actor_id":   notification.UserID,
		"entity_id":  notification.RelatedID,
		"is_read":    0,
		"created_at": notification.CreatedAt,
	})

	return notification, nil
}
