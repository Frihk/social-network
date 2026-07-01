const API_URL = "http://localhost:8080/api";

export async function getFeed() {
  const res = await fetch(`${API_URL}/posts`, { credentials: "include" });
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
  const res = await fetch(`${API_URL}/users/${userId}/posts`, { credentials: "include" });
  if (!res.ok) throw new Error("Failed to fetch user posts");
  return res.json();
}

export async function getGroupPosts(groupId) {
  const res = await fetch(`${API_URL}/groups/${groupId}/posts`, { credentials: "include" });
  if (!res.ok) throw new Error("Failed to fetch group posts");
  return res.json();
}

export async function updatePost(postId, formData) {
  const res = await fetch(`${API_URL}/posts/${postId}`, {
    method: "PUT",
    body: formData,
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to update post");
  return res.json();
}

export async function deletePost(postId) {
  const res = await fetch(`${API_URL}/posts/${postId}`, {
    method: "DELETE",
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to delete post");
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
  const res = await fetch(`${API_URL}/posts/${postId}/comments`, { credentials: "include" });
  if (!res.ok) throw new Error("Failed to fetch comments");
  return res.json();
}

export const reactToPost = async (postId, emoji) => {
  const response = await fetch(`http://localhost:8080/api/posts/${postId}/react`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      // Include your authorization token here if required by middleware
      "Authorization": `Bearer ${localStorage.getItem("token")}` 
    },
    body: JSON.stringify({ emoji }),
  });

  if (!response.ok) {
    throw new Error("Failed to react to post");
  }

  return response.json();
};