const API_URL = 'http://localhost:8080/api';

export async function getChatHistory(userId) {
  const response = await fetch(`${API_URL}/chat/${userId}`, {
    credentials: 'include',
  });

  if (!response.ok) {
    if (response.status === 403) {
      throw new Error('No follow relationship with this user');
    }
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch message history');
  }

  return response.json();
}

export async function getGroupMessages(groupId) {
  const response = await fetch(`${API_URL}/groups/${groupId}/messages`, {
    credentials: 'include',
  });

  if (!response.ok) {
    if (response.status === 403) {
      throw new Error('Not a member of this group');
    }
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch group messages');
  }

  return response.json();
}

export async function sendPrivateMessage(userId, content) {
  const response = await fetch(`${API_URL}/chat/${userId}/messages`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({ content }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to send message');
  }

  return response.json();
}

export async function sendGroupMessage(groupId, content) {
  const response = await fetch(`${API_URL}/groups/${groupId}/messages`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({ content }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to send group message');
  }

  return response.json();
}

export async function getDMEligibleUsers() {
  const response = await fetch(`${API_URL}/chat/users`, {
    credentials: 'include',
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch DM eligible users');
  }

  return response.json();
}
