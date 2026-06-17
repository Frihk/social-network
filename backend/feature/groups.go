package feature

import (
	"encoding/json"
	"net/http"
	"social/models"
	"social/pkg/handlers"
	"social/queries"
	"strings"
)

// GroupHandlers contains all group-related HTTP handlers
type GroupHandlers struct {
	queries *queries.GroupQueries
}

// NewGroupHandlers creates a new GroupHandlers instance
func NewGroupHandlers(q *queries.GroupQueries) *GroupHandlers {
	return &GroupHandlers{queries: q}
}

// extractIDFromPath extracts the last path segment as ID
func extractIDFromPath(path, prefix string) string {
	remainder := strings.TrimPrefix(path, prefix)
	parts := strings.Split(remainder, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// ListGroups returns all groups - GET /api/groups
func (h *GroupHandlers) ListGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.queries.GetAllGroups()
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if groups == nil {
		groups = []*models.Group{}
	}
	handlers.JSONResponse(w, groups, http.StatusOK)
}

// CreateGroup creates a new group - POST /api/groups
func (h *GroupHandlers) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req models.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(handlers.UserIDKey{}).(string)
	if !ok {
		handlers.JSONError(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	group, err := h.queries.CreateGroup(req.Title, req.Description, userID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handlers.JSONResponse(w, group, http.StatusCreated)
}

// GetGroupDetails returns group details - GET /api/groups/{id}
func (h *GroupHandlers) GetGroupDetails(w http.ResponseWriter, r *http.Request) {
	groupID := extractIDFromPath(r.URL.Path, "/api/groups/")

	// Get user ID from context
	userID, ok := r.Context().Value(handlers.UserIDKey{}).(string)
	if !ok {
		handlers.JSONError(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	// Check if user is an accepted member
	status, err := h.queries.GetMembershipStatus(groupID, userID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status != "accepted" {
		handlers.JSONError(w, "only accepted members can view full group details", http.StatusForbidden)
		return
	}

	group, err := h.queries.GetGroupByID(groupID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	members, err := h.queries.GetGroupMembers(groupID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	acceptedCount := 0
	for _, member := range members {
		if member.Status == "accepted" {
			acceptedCount++
		}
	}

	if members == nil {
		members = []*models.GroupMember{}
	}

	response := models.GroupDetailResponse{
		Group:         group,
		Members:       members,
		MemberCount:   len(members),
		AcceptedCount: acceptedCount,
	}

	handlers.JSONResponse(w, response, http.StatusOK)
}

// InviteUserToGroup invites a user to a group - POST /api/groups/{id}/invite
func (h *GroupHandlers) InviteUserToGroup(w http.ResponseWriter, r *http.Request) {
	groupID := extractIDFromPath(r.URL.Path, "/api/groups/")
	groupID = strings.TrimSuffix(groupID, "/invite")

	var req models.InviteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(handlers.UserIDKey{}).(string)
	if !ok {
		handlers.JSONError(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	// Check if requester is the group creator
	group, err := h.queries.GetGroupByID(groupID)
	if err != nil {
		handlers.JSONError(w, "group not found", http.StatusNotFound)
		return
	}

	if group.CreatorID != userID {
		handlers.JSONError(w, "only group creator can invite users", http.StatusForbidden)
		return
	}

	// Check if user is already a member
	status, err := h.queries.GetMembershipStatus(groupID, req.UserID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status != "" {
		handlers.JSONError(w, "user is already a member of this group", http.StatusBadRequest)
		return
	}

	// Create invitation
	member, err := h.queries.InviteUserToGroup(groupID, req.UserID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create notification for invited user
	_, err = h.queries.CreateNotification(
		req.UserID,
		"group_invite",
		groupID,
		"You have been invited to join the group: "+group.Title,
	)
	if err != nil {
		handlers.JSONError(w, "failed to create notification", http.StatusInternalServerError)
		return
	}

	handlers.JSONResponse(w, member, http.StatusCreated)
}

// RequestToJoinGroup creates a request to join a group - POST /api/groups/{id}/request
func (h *GroupHandlers) RequestToJoinGroup(w http.ResponseWriter, r *http.Request) {
	groupID := extractIDFromPath(r.URL.Path, "/api/groups/")
	groupID = strings.TrimSuffix(groupID, "/request")

	// Get user ID from context
	userID, ok := r.Context().Value(handlers.UserIDKey{}).(string)
	if !ok {
		handlers.JSONError(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	// Check if group exists
	group, err := h.queries.GetGroupByID(groupID)
	if err != nil {
		handlers.JSONError(w, "group not found", http.StatusNotFound)
		return
	}

	// Check if user is already a member
	status, err := h.queries.GetMembershipStatus(groupID, userID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status != "" {
		handlers.JSONError(w, "you are already a member of this group", http.StatusBadRequest)
		return
	}

	// Create request
	member, err := h.queries.RequestToJoinGroup(groupID, userID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create notification for group creator
	_, err = h.queries.CreateNotification(
		group.CreatorID,
		"group_request",
		groupID,
		"A user has requested to join your group: "+group.Title,
	)
	if err != nil {
		handlers.JSONError(w, "failed to create notification", http.StatusInternalServerError)
		return
	}

	handlers.JSONResponse(w, member, http.StatusCreated)
}

// AcceptMember accepts a pending member request - PUT /api/groups/{id}/members/{userId}/accept
func (h *GroupHandlers) AcceptMember(w http.ResponseWriter, r *http.Request) {
	// Extract groupID and userId from path like /api/groups/{id}/members/{userId}/accept
	parts := strings.Split(r.URL.Path, "/")
	var groupID, userID string
	for i, part := range parts {
		if part == "groups" && i+1 < len(parts) {
			groupID = parts[i+1]
		}
		if part == "members" && i+1 < len(parts) {
			userID = parts[i+1]
		}
	}

	// Get authenticated user ID from context
	authUserID, ok := r.Context().Value(handlers.UserIDKey{}).(string)
	if !ok {
		handlers.JSONError(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	// Check if requester is the group creator
	group, err := h.queries.GetGroupByID(groupID)
	if err != nil {
		handlers.JSONError(w, "group not found", http.StatusNotFound)
		return
	}

	if group.CreatorID != authUserID {
		handlers.JSONError(w, "only group creator can accept members", http.StatusForbidden)
		return
	}

	// Accept the membership
	if err := h.queries.AcceptMembership(groupID, userID); err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create notification for the user
	_, err = h.queries.CreateNotification(
		userID,
		"group_accepted",
		groupID,
		"Your request to join the group has been accepted: "+group.Title,
	)

	handlers.JSONResponse(w, map[string]string{"message": "member accepted"}, http.StatusOK)
}

// DeclineMember declines/removes a member request - PUT /api/groups/{id}/members/{userId}/decline
func (h *GroupHandlers) DeclineMember(w http.ResponseWriter, r *http.Request) {
	// Extract groupID and userId from path like /api/groups/{id}/members/{userId}/decline
	parts := strings.Split(r.URL.Path, "/")
	var groupID, userID string
	for i, part := range parts {
		if part == "groups" && i+1 < len(parts) {
			groupID = parts[i+1]
		}
		if part == "members" && i+1 < len(parts) {
			userID = parts[i+1]
		}
	}

	// Get authenticated user ID from context
	authUserID, ok := r.Context().Value(handlers.UserIDKey{}).(string)
	if !ok {
		handlers.JSONError(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	// Check if requester is the group creator
	group, err := h.queries.GetGroupByID(groupID)
	if err != nil {
		handlers.JSONError(w, "group not found", http.StatusNotFound)
		return
	}

	if group.CreatorID != authUserID {
		handlers.JSONError(w, "only group creator can decline members", http.StatusForbidden)
		return
	}

	// Decline the membership
	if err := h.queries.DeclineMembership(groupID, userID); err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create notification for the user
	_, err = h.queries.CreateNotification(
		userID,
		"group_declined",
		groupID,
		"Your request to join the group has been declined: "+group.Title,
	)

	handlers.JSONResponse(w, map[string]string{"message": "member declined"}, http.StatusOK)
}

// AcceptGroupInvite accepts a group invitation - PUT /api/group-invites/{id}/accept
func (h *GroupHandlers) AcceptGroupInvite(w http.ResponseWriter, r *http.Request) {
	groupID := extractIDFromPath(r.URL.Path, "/api/group-invites/")
	groupID = strings.TrimSuffix(groupID, "/accept")

	// Get authenticated user ID from context
	userID, ok := r.Context().Value(handlers.UserIDKey{}).(string)
	if !ok {
		handlers.JSONError(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	// Verify the group exists
	_, err := h.queries.GetGroupByID(groupID)
	if err != nil {
		handlers.JSONError(w, "group not found", http.StatusNotFound)
		return
	}

	// Check that user has an invited status
	status, err := h.queries.GetMembershipStatus(groupID, userID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status != "invited" {
		handlers.JSONError(w, "only invited users can accept invitations", http.StatusForbidden)
		return
	}

	// Accept the invite
	if err := h.queries.AcceptMembership(groupID, userID); err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handlers.JSONResponse(w, map[string]string{"message": "invite accepted"}, http.StatusOK)
}

// DeclineGroupInvite declines a group invitation - PUT /api/group-invites/{id}/decline
func (h *GroupHandlers) DeclineGroupInvite(w http.ResponseWriter, r *http.Request) {
	groupID := extractIDFromPath(r.URL.Path, "/api/group-invites/")
	groupID = strings.TrimSuffix(groupID, "/decline")

	// Get authenticated user ID from context
	userID, ok := r.Context().Value(handlers.UserIDKey{}).(string)
	if !ok {
		handlers.JSONError(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	// Verify the group exists
	_, err := h.queries.GetGroupByID(groupID)
	if err != nil {
		handlers.JSONError(w, "group not found", http.StatusNotFound)
		return
	}

	// Check that user has an invited status
	status, err := h.queries.GetMembershipStatus(groupID, userID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status != "invited" {
		handlers.JSONError(w, "only invited users can decline invitations", http.StatusForbidden)
		return
	}

	// Decline the invite
	if err := h.queries.DeclineMembership(groupID, userID); err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handlers.JSONResponse(w, map[string]string{"message": "invite declined"}, http.StatusOK)
}

// ListGroupEvents returns all events for a group - GET /api/groups/{id}/events
func (h *GroupHandlers) ListGroupEvents(w http.ResponseWriter, r *http.Request) {
	groupID := extractIDFromPath(r.URL.Path, "/api/groups/")
	groupID = strings.TrimSuffix(groupID, "/events")

	// Get authenticated user ID from context
	userID, ok := r.Context().Value(handlers.UserIDKey{}).(string)
	if !ok {
		handlers.JSONError(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	// Check if user is an accepted member
	status, err := h.queries.GetMembershipStatus(groupID, userID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status != "accepted" {
		handlers.JSONError(w, "only accepted members can view events", http.StatusForbidden)
		return
	}

	events, err := h.queries.GetGroupEvents(groupID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if events == nil {
		events = []*models.Event{}
	}
	handlers.JSONResponse(w, events, http.StatusOK)
}

// CreateEvent creates a new event - POST /api/groups/{id}/events
func (h *GroupHandlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	groupID := extractIDFromPath(r.URL.Path, "/api/groups/")
	groupID = strings.TrimSuffix(groupID, "/events")

	var req models.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get authenticated user ID from context
	userID, ok := r.Context().Value(handlers.UserIDKey{}).(string)
	if !ok {
		handlers.JSONError(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	// Check if group exists
	group, err := h.queries.GetGroupByID(groupID)
	if err != nil {
		handlers.JSONError(w, "group not found", http.StatusNotFound)
		return
	}

	// Check if user is an accepted member
	status, err := h.queries.GetMembershipStatus(groupID, userID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status != "accepted" {
		handlers.JSONError(w, "only accepted members can create events", http.StatusForbidden)
		return
	}

	// Create event
	event, err := h.queries.CreateEvent(groupID, userID, req.Title, req.Description, req.EventTime)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get all accepted members and notify them
	members, err := h.queries.GetAcceptedMembers(groupID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, member := range members {
		if member.UserID != userID { // Don't notify the creator
			_, err = h.queries.CreateNotification(
				member.UserID,
				"event_created",
				event.ID,
				"A new event has been created in "+group.Title+": "+event.Title,
			)
		}
	}

	handlers.JSONResponse(w, event, http.StatusCreated)
}

// RespondToEvent creates or updates a response to an event - POST /api/events/{id}/respond
func (h *GroupHandlers) RespondToEvent(w http.ResponseWriter, r *http.Request) {
	eventID := extractIDFromPath(r.URL.Path, "/api/events/")
	eventID = strings.TrimSuffix(eventID, "/respond")

	var req models.RespondToEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get authenticated user ID from context
	userID, ok := r.Context().Value(handlers.UserIDKey{}).(string)
	if !ok {
		handlers.JSONError(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	// Check if event exists
	event, err := h.queries.GetEventByID(eventID)
	if err != nil {
		handlers.JSONError(w, "event not found", http.StatusNotFound)
		return
	}

	// Check if user is an accepted member of the group
	status, err := h.queries.GetMembershipStatus(event.GroupID, userID)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status != "accepted" {
		handlers.JSONError(w, "only accepted group members can respond to events", http.StatusForbidden)
		return
	}

	// Create or update response
	response, err := h.queries.RespondToEvent(eventID, userID, req.Response)
	if err != nil {
		handlers.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handlers.JSONResponse(w, response, http.StatusOK)
}
