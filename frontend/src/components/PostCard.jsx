"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { getComments, createComment } from "@/lib/posts";
import { getUserProfile } from "@/lib/auth";
import { getFollowing, followUser, unfollowUser } from "@/lib/followers";
import { useAuth } from "@/context/AuthContext";

export default function PostCard({ post }) {
  const { user: currentUser } = useAuth();
  const [followState, setFollowState] = useState(null);
  const [isPrivate, setIsPrivate] = useState(false);
  const [loadingFollow, setLoadingFollow] = useState(false);

  const [showComments, setShowComments] = useState(false);
  const [comments, setComments] = useState([]);
  const [commentContent, setCommentContent] = useState("");
  const [commentImage, setCommentImage] = useState(null);
  const [loadingComments, setLoadingComments] = useState(false);
  const [submitting, setSubmitting] = useState(false);

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

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString(undefined, {
      month: "short", day: "numeric", hour: "2-digit", minute: "2-digit",
    });
  };



  const authorName = post.author_name || `User ${post.user_id}`;

  return (
    <div className="card">
      <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "12px" }}>
        <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
          <div style={{ width: "40px", height: "40px", borderRadius: "50%", background: "#ddd", overflow: "hidden", flexShrink: 0 }}>
            {post.author_avatar && (
              <img src={`http://localhost:8080/${post.author_avatar}`} alt="avatar" style={{ width: "100%", height: "100%", objectFit: "cover" }} />
            )}
          </div>
          <div>
            <Link href={`/profile/${post.user_id}`} style={{ fontWeight: "600", fontSize: "15px", textDecoration: "none", color: "inherit", cursor: "pointer" }}>
              {authorName}
            </Link>
            <div style={{ fontSize: "12px", color: "var(--text-secondary)" }}>{formatDate(post.created_at)}</div>
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
              minWidth: "90px"
            }}
          >
            {followState === "following" ? "Following" : followState === "pending" ? "Requested" : "Follow"}
          </button>
        )}
      </div>

      <div style={{ fontSize: "15px", marginBottom: "12px", whiteSpace: "pre-wrap" }}>{post.content}</div>

      {post.image_path && (
        <div style={{ margin: "0 -16px 12px", borderTop: "1px solid var(--border-color)", borderBottom: "1px solid var(--border-color)" }}>
          <img src={`http://localhost:8080/${post.image_path}`} alt="Post content" style={{ width: "100%", display: "block" }} />
        </div>
      )}

      <div style={{ borderTop: "1px solid var(--border-color)", paddingTop: "8px" }}>
        <button onClick={handleToggleComments} className="btn-secondary" style={{ width: "100%", background: "transparent", color: "var(--text-secondary)" }}>
          {showComments ? "Hide Comments" : `Comments${comments.length > 0 ? ` (${comments.length})` : ""}`}
        </button>
      </div>

      {showComments && (
        <div style={{ marginTop: "12px" }}>
          {loadingComments ? (
            <p style={{ textAlign: "center", fontSize: "13px", color: "var(--text-secondary)" }}>Loading...</p>
          ) : (
            <>
              {comments.length === 0 && (
                <p style={{ fontSize: "13px", color: "var(--text-secondary)", textAlign: "center", marginBottom: "8px" }}>No comments yet.</p>
              )}
              {comments.map((comment) => (
                <div key={comment.id} style={{ display: "flex", gap: "8px", marginBottom: "8px" }}>
                  <div style={{ width: "32px", height: "32px", borderRadius: "50%", background: "#ddd", flexShrink: 0, overflow: "hidden" }}>
                    {comment.author_avatar && (
                      <img src={`http://localhost:8080/${comment.author_avatar}`} alt="avatar" style={{ width: "100%", height: "100%", objectFit: "cover" }} />
                    )}
                  </div>
                  <div style={{ background: "#f0f2f5", padding: "8px 12px", borderRadius: "18px", fontSize: "13px" }}>
                    <Link href={`/profile/${comment.user_id}`} style={{ fontWeight: "600", textDecoration: "none", color: "inherit", cursor: "pointer" }}>
                      {comment.author_name || `User ${comment.user_id}`}
                    </Link>
                    <div>{comment.content}</div>
                    {comment.image_path && (
                      <img src={`http://localhost:8080/${comment.image_path}`} alt="Comment" style={{ maxWidth: "200px", borderRadius: "8px", marginTop: "5px" }} />
                    )}
                  </div>
                </div>
              ))}

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
                  <div style={{ display: "flex", justifyContent: "space-between", marginTop: "4px", alignItems: "center" }}>
                    <input
                      type="file"
                      id={`comment-img-${post.id}`}
                      hidden
                      accept="image/jpeg,image/png,image/gif"
                      onChange={(e) => setCommentImage(e.target.files?.[0] || null)}
                    />
                    <label htmlFor={`comment-img-${post.id}`} style={{ fontSize: "12px", color: "var(--primary-color)", cursor: "pointer" }}>
                      {commentImage ? "Image selected" : "Add Image"}
                    </label>
                    <button type="submit" className="btn-primary" style={{ padding: "4px 12px", fontSize: "13px" }} disabled={submitting}>
                      {submitting ? "Posting..." : "Post"}
                    </button>
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
