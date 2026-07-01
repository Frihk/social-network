"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/context/AuthContext";
import { createPost } from "@/lib/posts";
import { getFollowers } from "@/lib/followers";

export default function PostForm({ onPostCreated, groupId }) {
  const { user } = useAuth();
  const [content, setContent] = useState("");
  const [privacy, setPrivacy] = useState("public");
  const [image, setImage] = useState(null);
  const [followers, setFollowers] = useState([]);
  const [selectedViewers, setSelectedViewers] = useState([]);
  const [loadingFollowers, setLoadingFollowers] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (privacy === "private" && user) {
      const fetchFollowers = async () => {
        setLoadingFollowers(true);
        try {
          const data = await getFollowers(user.id);
          setFollowers(data || []);
        } catch (err) {
          console.error("Failed to fetch followers", err);
        } finally {
          setLoadingFollowers(false);
        }
      };
      fetchFollowers();
    }
  }, [privacy, user]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!content.trim() && !image) return;

    setError('');
    setSubmitting(true);

    const formData = new FormData();
    formData.append("content", content);
    formData.append("privacy", privacy);
    if (groupId) formData.append("group_id", groupId);
    if (image) formData.append("image", image);
    if (privacy === "private") {
      selectedViewers.forEach((id) => formData.append("allowed_viewers", id));
    }

    try {
      await createPost(formData);
      setContent("");
      setPrivacy("public");
      setImage(null);
      setSelectedViewers([]);
      if (onPostCreated) onPostCreated();
    } catch (err) {
      setError(err.message || "Failed to create post");
    } finally {
      setSubmitting(false);
    }
  };

  const handleViewerToggle = (id) => {
    setSelectedViewers((prev) =>
      prev.includes(id) ? prev.filter((v) => v !== id) : [...prev, id]
    );
  };

  return (
    <div className="card">
      <div style={{ display: "flex", gap: "10px", marginBottom: "12px" }}>
        <div style={{ width: "40px", height: "40px", borderRadius: "50%", background: "#ddd" }} />
        <div style={{ flex: 1 }}>
          <form onSubmit={handleSubmit}>
            <textarea
              className="input-field"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder={`What's on your mind, ${user?.first_name || "User"}?`}
              style={{ borderRadius: "12px", minHeight: "60px" }}
            />
            <div style={{ display: "flex", justifyContent: "space-between", marginTop: "12px", alignItems: "center" }}>
              <div style={{ display: "flex", gap: "8px" }}>
                <select
                  value={privacy}
                  onChange={(e) => setPrivacy(e.target.value)}
                  style={{ border: "none", background: "#f0f2f5", borderRadius: "6px", padding: "4px 8px", fontSize: "12px", fontWeight: "600" }}
                >
                  <option value="public">Public</option>
                  <option value="almost_private">Followers</option>
                  <option value="private">Specific</option>
                </select>
                <input
                  type="file"
                  id="post-image"
                  hidden
                  accept="image/jpeg,image/png,image/gif"
                  onChange={(e) => setImage(e.target.files?.[0] || null)}
                />
                <label htmlFor="post-image" style={{ background: "#f0f2f5", borderRadius: "6px", padding: "4px 8px", fontSize: "12px", fontWeight: "600", cursor: "pointer" }}>
                  {image ? "Image selected" : "Photo/Video"}
                </label>
              </div>
              <button type="submit" className="btn-primary" style={{ padding: "6px 24px" }} disabled={submitting}>{submitting ? "Posting..." : "Post"}</button>
            </div>

            {error && (
              <div style={{ color: "red", fontSize: "13px", marginTop: "8px" }}>{error}</div>
            )}
            {privacy === "private" && (
              <div style={{ marginTop: "12px", padding: "12px", border: "1px solid var(--border-color)", borderRadius: "8px", background: "#f9f9f9" }}>
                <div style={{ fontSize: "13px", fontWeight: "600", marginBottom: "8px" }}>Select Viewers</div>
                {loadingFollowers ? (
                  <p style={{ fontSize: "12px" }}>Loading followers...</p>
                ) : (
                  <select
                    multiple
                    value={selectedViewers}
                    onChange={(e) => {
                      const options = [...e.target.selectedOptions];
                      const values = options.map(option => option.value);
                      setSelectedViewers(values);
                    }}
                    style={{ width: "100%", padding: "8px", borderRadius: "6px", border: "1px solid var(--border-color)", background: "#fff", minHeight: "80px" }}
                  >
                    {followers.map(f => (
                      <option key={f.id} value={f.follower_id} style={{ padding: "4px" }}>
                        User {f.follower_id}
                      </option>
                    ))}
                    {followers.length === 0 && <option disabled value="">No followers found.</option>}
                  </select>
                )}
              </div>
            )}
          </form>
        </div>
      </div>
    </div>
  );
}
