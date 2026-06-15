"use client";

import { useState } from "react";
import { getComments, createComment } from "@/lib/posts";

export default function PostCard({ post }) {
  const [showComments, setShowComments] = useState(false);
  const [comments, setComments] = useState([]);
  const [commentContent, setCommentContent] = useState("");
  const [commentImage, setCommentImage] = useState(null);
  const [loadingComments, setLoadingComments] = useState(false);

  const handleToggleComments = async () => {
    if (!showComments) {
      setLoadingComments(true);
      try {
        const data = await getComments(post.id);
        setComments(data || []);
      } catch (err) {
        console.error(err);
      } finally {
        setLoadingComments(false);
      }
    }
    setShowComments(!showComments);
  };

  const handleCommentSubmit = async (e) => {
    e.preventDefault();
    if (!commentContent.trim() && !commentImage) return;

    const formData = new FormData();
    formData.append("content", commentContent);
    if (commentImage) {
      formData.append("image", commentImage);
    }

    try {
      const newComment = await createComment(post.id, formData);
      setComments([...comments, newComment]);
      setCommentContent("");
      setCommentImage(null);
    } catch (err) {
      console.error(err);
    }
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString(undefined, { 
      month: 'short', 
      day: 'numeric', 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };

  const getBadgeClass = (privacy) => {
    if (privacy === 'public') return 'badge-public';
    if (privacy === 'private') return 'badge-private';
    return 'badge-almost';
  };

  return (
    <div className="card">
      <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "12px" }}>
        <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
          <div style={{ width: "40px", height: "40px", borderRadius: "50%", background: "#ddd" }} />
          <div>
            <div style={{ fontWeight: "600", fontSize: "15px" }}>User {post.user_id}</div>
            <div style={{ fontSize: "12px", color: "var(--text-secondary)" }}>{formatDate(post.created_at)}</div>
          </div>
        </div>
        <span className={`badge ${getBadgeClass(post.privacy)}`}>
          {post.privacy.replace('_', ' ')}
        </span>
      </div>

      <div style={{ fontSize: "15px", marginBottom: "12px", whiteSpace: "pre-wrap" }}>
        {post.content}
      </div>

      {post.image_path && (
        <div style={{ margin: "0 -16px 12px", borderTop: "1px solid var(--border-color)", borderBottom: "1px solid var(--border-color)" }}>
          <img
            src={`http://localhost:8080/${post.image_path}`}
            alt="Post content"
            style={{ width: "100%", display: "block" }}
          />
        </div>
      )}

      <div style={{ borderTop: "1px solid var(--border-color)", paddingTop: "8px" }}>
        <button 
          onClick={handleToggleComments} 
          className="btn-secondary" 
          style={{ width: "100%", background: "transparent", color: "var(--text-secondary)" }}
        >
          {comments.length > 0 ? `${comments.length} Comments` : "Comment"}
        </button>
      </div>

      {showComments && (
        <div style={{ marginTop: "12px" }}>
          {loadingComments ? (
            <p style={{ textAlign: "center", fontSize: "13px", color: "var(--text-secondary)" }}>Loading...</p>
          ) : (
            <>
              {comments.map((comment) => (
                <div key={comment.id} style={{ display: "flex", gap: "8px", marginBottom: "8px" }}>
                  <div style={{ width: "32px", height: "32px", borderRadius: "50%", background: "#ddd", flexShrink: 0 }} />
                  <div style={{ background: "#f0f2f5", padding: "8px 12px", borderRadius: "18px", fontSize: "13px" }}>
                    <div style={{ fontWeight: "600" }}>User {comment.user_id}</div>
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
                      onChange={(e) => setCommentImage(e.target.files[0])}
                    />
                    <label htmlFor={`comment-img-${post.id}`} style={{ fontSize: "12px", color: "var(--primary-color)", cursor: "pointer" }}>
                      {commentImage ? "Image selected" : "Add Image"}
                    </label>
                    <button type="submit" className="btn-primary" style={{ padding: "4px 12px", fontSize: "13px" }}>
                      Post
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
