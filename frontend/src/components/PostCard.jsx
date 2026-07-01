"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { getComments, createComment, updatePost, deletePost } from "@/lib/posts";
import { getUserProfile } from "@/lib/auth";
import { getFollowing, followUser, unfollowUser } from "@/lib/followers";
import { useAuth } from "@/context/AuthContext";
import { reactToPost } from "@/lib/posts"; 
import EmojiPicker from "./EmojiPicker";

export default function PostCard({ post }) {
  const { user: currentUser } = useAuth();
  const [followState, setFollowState] = useState(null);
  const [isPrivate, setIsPrivate] = useState(false);
  const [loadingFollow, setLoadingFollow] = useState(false);
  const [editing, setEditing] = useState(false);
  const [editContent, setEditContent] = useState("");
  const [editPrivacy, setEditPrivacy] = useState("public");

  const [showComments, setShowComments] = useState(false);
  const [comments, setComments] = useState([]);
  const [commentCount, setCommentCount] = useState(post.comment_count ?? 0);
  const [commentContent, setCommentContent] = useState("");
  const [commentImage, setCommentImage] = useState(null);
  const [loadingComments, setLoadingComments] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const [showPicker, setShowPicker] = useState(false);
  const [reactions, setReactions] = useState(post.reactions || []);

  const handleReact = async (emoji) => {
    try {
      await reactToPost(post.id, emoji);
      const existingReactionIndex = reactions.findIndex(r => r.user_id === currentUser?.id);
      let updatedReactions = [...reactions];
      if (existingReactionIndex !== -1) {
        if (updatedReactions[existingReactionIndex].emoji === emoji) {
          updatedReactions.splice(existingReactionIndex, 1);
        } else {
          updatedReactions[existingReactionIndex].emoji = emoji;
        }
      } else {
        updatedReactions.push({ user_id: currentUser?.id, emoji: emoji });
      }
      setReactions(updatedReactions);
      setShowPicker(false);
    } catch (error) {
      console.error("Error reacting to post:", error);
    }
  };

  const reactionCounts = reactions.reduce((acc, curr) => {
    acc[curr.emoji] = (acc[curr.emoji] || 0) + 1;
    return acc;
  }, {});

  useEffect(() => {
    if (currentUser && currentUser.id !== post.user_id) {
      const checkFollow = async () => {
        try {
          const profile = await getUserProfile(post.user_id);
          setIsPrivate(profile.is_private);
          
          const myFollowing = await getFollowing(currentUser.id);
          const followRel = myFollowing?.find(f => f.following_id === post.user_id);
          if (followRel) {
            setFollowState(followRel.status === "pending" ? "pending" : "following");
          } else {
            setFollowState("not_following");
          }
        } catch (e) {
          console.error(e);
        }
      };
      checkFollow();
    }
  }, [currentUser, post.user_id]);

  const handleFollowToggle = async () => {
    setLoadingFollow(true);
    try {
      if (followState === "following") {
        if (!window.confirm("Are you sure?")) {
          setLoadingFollow(false);
          return;
        }
        setFollowState("not_following");
        await unfollowUser(post.user_id);
      } else if (followState === "not_following") {
        setFollowState(isPrivate ? "pending" : "following");
        await followUser(post.user_id);
      }
    } catch (err) {
      console.error(err);
      // Revert state if needed in production
    } finally {
      setLoadingFollow(false);
    }
  };

  const fetchComments = async () => {
    setLoadingComments(true);
    try {
      const data = await getComments(post.id);
      setComments(data || []);
      setCommentCount(data?.length ?? 0);
    } catch (err) {
      console.error(err);
    } finally {
      setLoadingComments(false);
    }
  };

  const handleToggleComments = async () => {
    if (!showComments) await fetchComments();
    setShowComments((prev) => !prev);
  };

  const handleCommentSubmit = async (e) => {
    e.preventDefault();
    if (!commentContent.trim() && !commentImage) return;

    const formData = new FormData();
    formData.append("content", commentContent);
    if (commentImage) formData.append("image", commentImage);

    setSubmitting(true);
    try {
      await createComment(post.id, formData);
      setCommentContent("");
      setCommentImage(null);
      await fetchComments();
    } catch (err) {
      console.error(err);
    } finally {
      setSubmitting(false);
    }
  };

  const handleEdit = () => {
    setEditContent(post.content);
    setEditPrivacy(post.privacy || "public");
    setEditing(true);
  };

  const handleEditSubmit = async (e) => {
    e.preventDefault();
    if (!editContent.trim()) return;
    try {
      const formData = new FormData();
      formData.append("content", editContent);
      formData.append("privacy", editPrivacy);
      await updatePost(post.id, formData);
      setEditing(false);
      window.location.reload();
    } catch (err) {
      console.error("Failed to update post:", err);
    }
  };

  const handleDelete = async () => {
    if (!window.confirm("Delete this post?")) return;
    try {
      await deletePost(post.id);
      window.location.reload();
    } catch (err) {
      console.error("Failed to delete post:", err);
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString(undefined, {
      month: "short", day: "numeric", hour: "2-digit", minute: "2-digit",
    });
  };



  const authorName = post.author_name || `User ${post.user_id}`;

  return (
  <div className="card">

    {/* Header: avatar + author + follow button */}
    <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "12px" }}>
      <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
        <div style={{ width: "40px", height: "40px", borderRadius: "50%", background: "#ddd", overflow: "hidden", flexShrink: 0 }}>
          {post.author_avatar && (
            <img
              src={`http://localhost:8080/${post.author_avatar}`}
              alt="avatar"
              style={{ width: "100%", height: "100%", objectFit: "cover" }}
            />
          )}
        </div>
        <div>
          <Link
            href={`/profile/${post.user_id}`}
            style={{ fontWeight: "600", fontSize: "15px", textDecoration: "none", color: "inherit" }}
          >
            {authorName}
          </Link>
          <div style={{ fontSize: "12px", color: "var(--text-secondary)" }}>
            {formatDate(post.created_at)}
          </div>
        </div>
      </div>

      {currentUser && currentUser.id !== post.user_id && followState && (
        <button
          onClick={handleFollowToggle}
          disabled={followState === "pending" || loadingFollow}
          className={followState === "following" || followState === "pending" ? "btn-secondary" : "btn-primary"}
          style={{
            padding: "4px 12px",
            fontSize: "13px",
            height: "fit-content",
            opacity: followState === "pending" ? 0.6 : 1,
            cursor: followState === "pending" ? "not-allowed" : "pointer",
            minWidth: "90px",
          }}
        >
          {followState === "following" ? "Following" : followState === "pending" ? "Requested" : "Follow"}
        </button>
      )}
      {currentUser && currentUser.id === post.user_id && !editing && (
        <div style={{ display: "flex", gap: "4px" }}>
          <button onClick={handleEdit} style={{ background: "none", border: "1px solid #e5e7eb", borderRadius: "6px", padding: "2px 8px", fontSize: "12px", cursor: "pointer", color: "#374151" }}>Edit</button>
          <button onClick={handleDelete} style={{ background: "none", border: "1px solid #fecaca", borderRadius: "6px", padding: "2px 8px", fontSize: "12px", cursor: "pointer", color: "#dc2626" }}>Delete</button>
        </div>
      )}
    </div>

    {/* Post body */}
    {editing ? (
      <form onSubmit={handleEditSubmit} style={{ marginBottom: "12px" }}>
        <textarea
          value={editContent}
          onChange={(e) => setEditContent(e.target.value)}
          style={{ width: "100%", padding: "8px 12px", borderRadius: "8px", border: "1px solid #d1d5db", fontSize: "15px", minHeight: "60px", fontFamily: "inherit" }}
          required
        />
        <div style={{ display: "flex", gap: "8px", marginTop: "8px", alignItems: "center" }}>
          <select
            value={editPrivacy}
            onChange={(e) => setEditPrivacy(e.target.value)}
            style={{ border: "1px solid #d1d5db", borderRadius: "6px", padding: "4px 8px", fontSize: "12px" }}
          >
            <option value="public">Public</option>
            <option value="almost_private">Followers</option>
            <option value="private">Specific</option>
          </select>
          <button type="submit" style={{ background: "#3b82f6", color: "white", border: "none", padding: "6px 16px", borderRadius: "6px", fontSize: "13px", cursor: "pointer" }}>Save</button>
          <button type="button" onClick={() => setEditing(false)} style={{ background: "none", border: "1px solid #d1d5db", padding: "6px 16px", borderRadius: "6px", fontSize: "13px", cursor: "pointer" }}>Cancel</button>
        </div>
      </form>
    ) : (
      <div style={{ fontSize: "15px", marginBottom: "12px", whiteSpace: "pre-wrap" }}>
        {post.content}
      </div>
    )}

    {/* Post image */}
    {post.image_path && (
      <div style={{ margin: "0 -16px 12px", borderTop: "1px solid var(--border-color)", borderBottom: "1px solid var(--border-color)" }}>
        <img
          src={`http://localhost:8080/${post.image_path}`}
          alt="Post content"
          style={{ width: "100%", display: "block" }}
        />
      </div>
    )}

    {/* Action bar: React + Comments in one row */}
    <div style={{ display: "flex", gap: "8px", borderTop: "1px solid var(--border-color)", paddingTop: "8px", marginTop: "4px" }}>

      {/* React button + picker */}
      <div style={{ position: "relative" }}>
        <button
          onClick={() => setShowPicker(!showPicker)}
          className="btn-secondary"
          style={{ background: "transparent", color: "var(--text-secondary)" }}
        >
          🎃React🎃
        </button>

        {showPicker && (
          <div style={{
            position: "absolute",
            bottom: "36px",
            left: "0",
            zIndex: 10,
            background: "white",
            boxShadow: "0 4px 6px rgba(0,0,0,0.1)",
            borderRadius: "8px",
            border: "1px solid #e5e7eb",
          }}>
            <EmojiPicker onSelect={handleReact} />
          </div>
        )}
      </div>

      {/* Comments toggle */}
      <button
        onClick={handleToggleComments}
        className="btn-secondary"
        style={{ flex: 1, background: "transparent", color: "var(--text-secondary)" }}
      >
        💬 {showComments ? "Hide comments" : `Comments (${commentCount})`}
      </button>
    </div>

    {/* Reaction pills — only rendered when there are reactions */}
    {Object.keys(reactionCounts).length > 0 && (
      <div style={{ display: "flex", flexWrap: "wrap", gap: "6px", marginTop: "8px" }}>
        {Object.entries(reactionCounts).map(([emoji, count]) => (
          <button
            key={emoji}
            onClick={() => handleReact(emoji)}
            style={{
              display: "flex",
              alignItems: "center",
              gap: "4px",
              background: reactions.some(r => r.user_id === currentUser?.id && r.emoji === emoji)
                ? "#dbeafe"
                : "#f3f4f6",
              padding: "2px 8px",
              borderRadius: "9999px",
              fontSize: "13px",
              border: "1px solid #e5e7eb",
              cursor: "pointer",
            }}
          >
            {emoji}
            <span style={{ color: "#4b5563" }}>{count}</span>
          </button>
        ))}
      </div>
    )}

    {/* Comments section */}
    {showComments && (
      <div style={{ marginTop: "12px" }}>
        {loadingComments ? (
          <p style={{ textAlign: "center", fontSize: "13px", color: "var(--text-secondary)" }}>Loading...</p>
        ) : (
          <>
            {comments.length === 0 && (
              <p style={{ fontSize: "13px", color: "var(--text-secondary)", textAlign: "center", marginBottom: "8px" }}>
                No comments yet.
              </p>
            )}

            {comments.map((comment) => (
              <div key={comment.id} style={{ display: "flex", gap: "8px", marginBottom: "8px" }}>
                <div style={{ width: "32px", height: "32px", borderRadius: "50%", background: "#ddd", flexShrink: 0, overflow: "hidden" }}>
                  {comment.author_avatar && (
                    <img
                      src={`http://localhost:8080/${comment.author_avatar}`}
                      alt="avatar"
                      style={{ width: "100%", height: "100%", objectFit: "cover" }}
                    />
                  )}
                </div>
                <div style={{ background: "#f0f2f5", padding: "8px 12px", borderRadius: "18px", fontSize: "13px" }}>
                  <Link
                    href={`/profile/${comment.user_id}`}
                    style={{ fontWeight: "600", textDecoration: "none", color: "inherit" }}
                  >
                    {comment.author_name || `User ${comment.user_id}`}
                  </Link>
                  <div>{comment.content}</div>
                  {comment.image_path && (
                    <img
                      src={`http://localhost:8080/${comment.image_path}`}
                      alt="Comment"
                      style={{ maxWidth: "200px", borderRadius: "8px", marginTop: "5px" }}
                    />
                  )}
                </div>
              </div>
            ))}

            {/* Comment form */}
            <form onSubmit={handleCommentSubmit} style={{ display: "flex", gap: "8px", marginTop: "12px", alignItems: "flex-start" }}>
              <div style={{ width: "32px", height: "32px", borderRadius: "50%", background: "#ddd", flexShrink: 0 }} />
              <div style={{ flex: 1 }}>
                <textarea
                  className="input-field"
                  value={commentContent}
                  onChange={(e) => setCommentContent(e.target.value)}
                  placeholder="Write a comment..."
                  rows="1"
                  style={{ padding: "8px 12px", borderRadius: "20px" }}
                />

                <div style={{ display: "flex", flexDirection: "column", gap: "8px", marginTop: "4px" }}>
                  {commentImage && (
                    <div style={{ position: "relative", display: "inline-block", width: "fit-content" }}>
                      <img
                        src={URL.createObjectURL(commentImage)}
                        alt="Preview"
                        style={{ maxWidth: "150px", borderRadius: "8px" }}
                      />
                      <button
                        type="button"
                        onClick={() => setCommentImage(null)}
                        style={{ position: "absolute", top: "-5px", right: "-5px", background: "red", color: "white", borderRadius: "50%", border: "none", width: "20px", height: "20px", cursor: "pointer", display: "flex", alignItems: "center", justifyContent: "center", fontSize: "10px" }}
                      >
                        X
                      </button>
                    </div>
                  )}

                  <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <input
                      type="file"
                      id={`comment-img-${post.id}`}
                      hidden
                      accept="image/jpeg,image/png,image/gif"
                      onChange={(e) => setCommentImage(e.target.files?.[0] || null)}
                    />
                    <label htmlFor={`comment-img-${post.id}`} style={{ fontSize: "12px", color: "var(--primary-color)", cursor: "pointer" }}>
                      {commentImage ? "Change image" : "Add image"}
                    </label>
                    <button
                      type="submit"
                      className="btn-primary"
                      style={{ padding: "4px 12px", fontSize: "13px" }}
                      disabled={submitting}
                    >
                      {submitting ? "Posting..." : "Post"}
                    </button>
                  </div>
                </div>
              </div>
            </form>
          </>
        )}
      </div>
    )}

  </div>
);
}

