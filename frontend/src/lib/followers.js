const API_URL = "http://localhost:8080/api";

export async function followUser(userId) {
  const res = await fetch(`${API_URL}/users/${userId}/follow`, {
    method: "POST",
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to follow user");
  return res.json();
}

export async function acceptFollow(requestId) {
  const res = await fetch(`${API_URL}/followers/${requestId}/accept`, {
    method: "PUT",
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to accept follow request");
  return res.json();
}

export async function declineFollow(requestId) {
  const res = await fetch(`${API_URL}/followers/${requestId}/decline`, {
    method: "PUT",
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to decline follow request");
  return res.json();
}

export async function unfollowUser(userId) {
  const res = await fetch(`${API_URL}/users/${userId}/unfollow`, {
    method: "POST",
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to unfollow user");
  return res.json();
}

export async function getFollowers(userId) {
  const res = await fetch(`${API_URL}/users/${userId}/followers`, {
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to fetch followers");
  return res.json();
}

export async function getFollowing(userId) {
  const res = await fetch(`${API_URL}/users/${userId}/following`, {
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to fetch following");
  return res.json();
}
