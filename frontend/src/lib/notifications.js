const API_URL = "http://localhost:8080/api";

export async function fetchNotifications() {
  const response = await fetch(`${API_URL}/notifications`, {
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch notifications");
  }

  const payload = await response.json();
  return payload.data || [];
}

export async function markNotificationRead(notificationId) {
  const response = await fetch(`${API_URL}/notifications/${notificationId}/read`, {
    method: "PUT",
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Failed to mark notification as read");
  }

  return response.json();
}

export async function markAllNotificationsRead() {
  const response = await fetch(`${API_URL}/notifications/read-all`, {
    method: "PUT",
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Failed to mark notifications as read");
  }

  return response.json();
}
