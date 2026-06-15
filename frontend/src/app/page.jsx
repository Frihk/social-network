"use client";

import { useState, useEffect, useCallback } from "react";
import PostForm from "@/components/PostForm";
import PostCard from "@/components/PostCard";
import { getFeed } from "@/lib/posts";

export default function Home() {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchPosts = useCallback(async () => {
    setLoading(true);
    try {
      const data = await getFeed();
      setPosts(data || []);
      setError(null);
    } catch (err) {
      console.error(err);
      setError("Failed to load feed. Please try again later.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchPosts();
  }, [fetchPosts]);

  return (
    <main className="container">
      <div style={{ marginBottom: "20px" }}>
        <h1 style={{ fontSize: "24px", fontWeight: "700" }}>Feed</h1>
      </div>
      
      <PostForm onPostCreated={fetchPosts} />
      
      {loading ? (
        <div style={{ textAlign: "center", marginTop: "40px" }}>
          <p style={{ color: "var(--text-secondary)" }}>Loading your feed...</p>
        </div>
      ) : error ? (
        <div className="card" style={{ textAlign: "center", color: "#d32f2f" }}>
          {error}
        </div>
      ) : posts.length === 0 ? (
        <div className="card" style={{ textAlign: "center", padding: "40px" }}>
          <h2 style={{ fontSize: "18px", marginBottom: "8px" }}>Your feed is empty</h2>
          <p style={{ color: "var(--text-secondary)" }}>Follow some people or create your first post!</p>
        </div>
      ) : (
        <div>
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </div>
      )}
    </main>
  );
}
