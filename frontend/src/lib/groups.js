const API_URL = 'http://localhost:8080/api';

export async function getUserGroups() {
  const response = await fetch(`${API_URL}/user/groups`, {
    credentials: 'include',
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch groups');
  }

  return response.json();
}