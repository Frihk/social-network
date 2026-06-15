const API_URL = "http://localhost:8080/api";

export async function getFeed() {
  const res = await fetch(`${API_URL}/posts`, {
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to fetch feed");
  return res.json();
}

export async function createPost(formData) {
  const res = await fetch(`${API_URL}/posts`, {
    method: "POST",
    body: formData,
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to create post");
  return res.json();
}

export async function getUserPosts(userId) {
  const res = await fetch(`${API_URL}/users/${userId}/posts`, {
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to fetch user posts");
  return res.json();
}

export async function createComment(postId, formData) {
  const res = await fetch(`${API_URL}/posts/${postId}/comments`, {
    method: "POST",
    body: formData,
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to create comment");
  return res.json();
}

export async function getComments(postId) {
  const res = await fetch(`${API_URL}/posts/${postId}/comments`, {
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to fetch comments");
  return res.json();
}
