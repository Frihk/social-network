const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

/**
 * Helper to get the auth headers.
 * In a real app this would use a proper token. For this backend, it uses X-User-ID.
 */
function getHeaders(userId) {
  const headers = {
    'Content-Type': 'application/json',
  };
  if (userId) {
    headers['X-User-ID'] = userId;
  }
  return headers;
}

// ==========================================
// GROUPS
// ==========================================

export async function fetchGroups(userId) {
  const response = await fetch(`${API_URL}/groups`, {
    headers: getHeaders(userId),
    cache: 'no-store'
  });
  if (!response.ok) throw new Error('Failed to fetch groups');
  return response.json();
}

export async function createGroup(userId, data) {
  const response = await fetch(`${API_URL}/groups`, {
    method: 'POST',
    headers: getHeaders(userId),
    body: JSON.stringify(data),
  });
  if (!response.ok) throw new Error('Failed to create group');
  return response.json();
}

export async function fetchGroupDetails(userId, groupId) {
  const response = await fetch(`${API_URL}/groups/${groupId}`, {
    headers: getHeaders(userId),
    cache: 'no-store'
  });
  if (!response.ok) throw new Error('Failed to fetch group details');
  return response.json();
}

// ==========================================
// GROUP MEMBERSHIP
// ==========================================

export async function inviteUserToGroup(userId, groupId, inviteeId) {
  const response = await fetch(`${API_URL}/groups/${groupId}/invite`, {
    method: 'POST',
    headers: getHeaders(userId),
    body: JSON.stringify({ user_id: inviteeId }),
  });
  if (!response.ok) throw new Error('Failed to invite user');
  return response.json();
}

export async function requestToJoinGroup(userId, groupId) {
  const response = await fetch(`${API_URL}/groups/${groupId}/request`, {
    method: 'POST',
    headers: getHeaders(userId),
  });
  if (!response.ok) throw new Error('Failed to request to join group');
  return response.json();
}

export async function acceptMemberRequest(userId, groupId, memberId) {
  const response = await fetch(`${API_URL}/groups/${groupId}/members/${memberId}/accept`, {
    method: 'PUT',
    headers: getHeaders(userId),
  });
  if (!response.ok) throw new Error('Failed to accept member request');
  return response.json();
}

export async function declineMemberRequest(userId, groupId, memberId) {
  const response = await fetch(`${API_URL}/groups/${groupId}/members/${memberId}/decline`, {
    method: 'PUT',
    headers: getHeaders(userId),
  });
  if (!response.ok) throw new Error('Failed to decline member request');
  return response.json();
}

export async function acceptGroupInvitation(userId, inviteId) {
  const response = await fetch(`${API_URL}/group-invites/${inviteId}/accept`, {
    method: 'PUT',
    headers: getHeaders(userId),
  });
  if (!response.ok) throw new Error('Failed to accept group invitation');
  return response.json();
}

export async function declineGroupInvitation(userId, inviteId) {
  const response = await fetch(`${API_URL}/group-invites/${inviteId}/decline`, {
    method: 'PUT',
    headers: getHeaders(userId),
  });
  if (!response.ok) throw new Error('Failed to decline group invitation');
  return response.json();
}

// ==========================================
// EVENTS
// ==========================================

export async function fetchGroupEvents(userId, groupId) {
  const response = await fetch(`${API_URL}/groups/${groupId}/events`, {
    headers: getHeaders(userId),
    cache: 'no-store'
  });
  if (!response.ok) throw new Error('Failed to fetch events');
  return response.json();
}

export async function createEvent(userId, groupId, data) {
  const response = await fetch(`${API_URL}/groups/${groupId}/events`, {
    method: 'POST',
    headers: getHeaders(userId),
    body: JSON.stringify(data),
  });
  if (!response.ok) throw new Error('Failed to create event');
  return response.json();
}

export async function respondToEvent(userId, eventId, responseType) {
  // responseType should be 'going' or 'not_going'
  const response = await fetch(`${API_URL}/events/${eventId}/respond`, {
    method: 'POST',
    headers: getHeaders(userId),
    body: JSON.stringify({ response: responseType }),
  });
  if (!response.ok) throw new Error('Failed to respond to event');
  return response.json();
}