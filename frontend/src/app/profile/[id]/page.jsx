"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/context/AuthContext";
import { getUserProfile, updateProfilePrivacy, updateUserProfile } from "@/lib/auth";
import { getUserPosts } from "@/lib/posts";
import { getFollowers, getFollowing, followUser, unfollowUser, acceptFollow, declineFollow } from "@/lib/followers";
import PostCard from "@/components/PostCard";

export default function ProfilePage() {
  const { id: profileId } = useParams();
  const { user: currentUser } = useAuth();
  const isOwnProfile = currentUser && currentUser.id === profileId;

  const [profileUser, setProfileUser] = useState(null);
  const [posts, setPosts] = useState([]);
  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const [pendingRequests, setPendingRequests] = useState([]);
  const [followState, setFollowState] = useState("not_following");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [updatingPrivacy, setUpdatingPrivacy] = useState(false);
  const [updatingAvatar, setUpdatingAvatar] = useState(false);

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const [userData, postsData, followersData, followingData] = await Promise.all([
        getUserProfile(profileId),
        getUserPosts(profileId),
        getFollowers(profileId),
        getFollowing(profileId),
      ]);

      setProfileUser(userData);
      setPosts(postsData || []);
      setFollowers(followersData?.filter((f) => f.status === "accepted") || []);
      setFollowing(followingData || []);

      if (isOwnProfile) {
        setPendingRequests(followersData?.filter((f) => f.status === "pending") || []);
      }

      if (currentUser && !isOwnProfile) {
        const myFollowing = await getFollowing(currentUser.id);
        const followRel = myFollowing?.find((f) => f.following_id === profileId);
        if (followRel) {
          setFollowState(followRel.status === "pending" ? "pending" : "following");
        } else {
          setFollowState("not_following");
        }
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [profileId, currentUser, isOwnProfile]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleFollowToggle = async () => {
    try {
      if (followState === "following") {
        if (!window.confirm("Are you sure?")) return;
        setFollowState("not_following");
        await unfollowUser(profileId);
      } else if (followState === "not_following") {
        setFollowState(profileUser?.is_private ? "pending" : "following");
        await followUser(profileId);
      }
      fetchData();
    } catch (err) {
      setError(err.message);
      fetchData();
    }
  };

  const handleRequest = async (requestId, action) => {
    try {
      if (action === "accept") await acceptFollow(requestId);
      else await declineFollow(requestId);
      fetchData();
    } catch (err) {
      setError(err.message);
    }
  };

  const handlePrivacyToggle = async () => {
    if (!isOwnProfile) return;
    if (!window.confirm("Are you sure?")) return;
    setUpdatingPrivacy(true);
    try {
      await updateProfilePrivacy(profileId, !profileUser.is_private);
      setProfileUser((prev) => ({ ...prev, is_private: !prev.is_private }));
    } catch (err) {
      setError(err.message);
    } finally {
      setUpdatingPrivacy(false);
    }
  };

  const handleAvatarChange = async (e) => {
    const avatar = e.target.files?.[0];
    if (!avatar || !isOwnProfile) return;

    setUpdatingAvatar(true);
    setError("");
    try {
      const formData = new FormData();
      formData.append("avatar", avatar);
      const updatedUser = await updateUserProfile(profileId, formData);
      setProfileUser(updatedUser);
    } catch (err) {
      setError(err.message);
    } finally {
      setUpdatingAvatar(false);
      e.target.value = "";
    }
  };

  if (loading) {
    return <div className="container" style={{ textAlign: "center" }}><p>Loading profile...</p></div>;
  }

  if (error) {
    return <div className="container" style={{ textAlign: "center" }}><p style={{ color: "#dc2626" }}>{error}</p></div>;
  }

  return (
    <div className="container">
      <div style={{ marginBottom: "16px", marginTop: "16px" }}>
        <Link href="/" style={{ color: "var(--primary-color)", textDecoration: "none", fontWeight: "600" }}>&larr; Back to Feed</Link>
      </div>
      <div className="card" style={{ padding: "40px", textAlign: "center", marginBottom: "24px" }}>
        <div style={{
          width: "120px", height: "120px", borderRadius: "50%", background: "#ddd",
          margin: "0 auto 16px", overflow: "hidden",
          display: "flex", alignItems: "center", justifyContent: "center",
        }}>
          {profileUser?.avatar ? (
            <img src={profileUser.avatar} alt="avatar" style={{ width: "100%", height: "100%", objectFit: "cover" }} />
          ) : (
            <span style={{ fontSize: "32px", color: "#6b7280" }}>
              {profileUser?.first_name?.[0]}{profileUser?.last_name?.[0]}
            </span>
          )}
        </div>
        <h1 style={{ fontSize: "28px", fontWeight: "700" }}>
          {profileUser ? `${profileUser.first_name} ${profileUser.last_name}` : "Unknown User"}
        </h1>
        {profileUser?.nickname && <p style={{ color: "var(--text-secondary)", marginTop: "4px" }}>@{profileUser.nickname}</p>}
        {profileUser?.about_me && <p style={{ marginTop: "8px", fontSize: "14px" }}>{profileUser.about_me}</p>}

        <div style={{ display: "flex", justifyContent: "center", gap: "24px", marginTop: "16px", marginBottom: "16px" }}>
          <div><span style={{ fontWeight: "700" }}>{posts.length}</span> <span style={{ color: "var(--text-secondary)", fontSize: "14px" }}>Posts</span></div>
          <div><span style={{ fontWeight: "700" }}>{followers.length}</span> <span style={{ color: "var(--text-secondary)", fontSize: "14px" }}>Followers</span></div>
          <div><span style={{ fontWeight: "700" }}>{following.length}</span> <span style={{ color: "var(--text-secondary)", fontSize: "14px" }}>Following</span></div>
        </div>

        {isOwnProfile && (
          <div style={{ marginTop: "16px" }}>
            <label className="btn-secondary" style={{ display: "inline-block", marginRight: "8px", minWidth: "140px", cursor: "pointer" }}>
              {updatingAvatar ? "Uploading..." : "Change Avatar"}
              <input
                type="file"
                accept="image/jpeg,image/png,image/gif"
                hidden
                disabled={updatingAvatar}
                onChange={handleAvatarChange}
              />
            </label>
            <button
              onClick={handlePrivacyToggle}
              disabled={updatingPrivacy}
              className={profileUser?.is_private ? "btn-secondary" : "btn-primary"}
              style={{ minWidth: "140px" }}
            >
              {updatingPrivacy
                ? "Updating..."
                : profileUser?.is_private
                ? "🔒 Private Profile"
                : "🌍 Public Profile"}
            </button>
          </div>
        )}

        {!isOwnProfile && currentUser && (
          <div style={{ marginTop: "16px" }}>
            <button 
              onClick={handleFollowToggle} 
              disabled={followState === "pending"}
              className={followState === "following" || followState === "pending" ? "btn-secondary" : "btn-primary"} 
              style={{ minWidth: "120px", opacity: followState === "pending" ? 0.6 : 1, cursor: followState === "pending" ? "not-allowed" : "pointer" }}
            >
              {followState === "following" ? "Following" : followState === "pending" ? "Requested" : "Follow"}
            </button>
          </div>
        )}
      </div>

      {!isOwnProfile && profileUser?.is_private && followState !== "following" ? (
        <div className="card" style={{ textAlign: "center", padding: "48px 20px", marginBottom: "24px" }}>
          <p style={{ fontSize: "18px", fontWeight: "600", marginBottom: "8px" }}>🔒 This account is private.</p>
          <p style={{ color: "var(--text-secondary)", fontSize: "15px" }}>Follow to see their posts and activity.</p>
        </div>
      ) : (
        <>
          {isOwnProfile && profileUser?.email && (
            <div className="card" style={{ marginBottom: "24px" }}>
              <h2 style={{ fontSize: "16px", fontWeight: "600", marginBottom: "12px" }}>Private Information</h2>
              <div style={{ fontSize: "14px" }}>
                <div><span style={{ fontWeight: "600" }}>Email:</span> {profileUser.email}</div>
                {profileUser.date_of_birth && (
                  <div><span style={{ fontWeight: "600" }}>Date of Birth:</span> {profileUser.date_of_birth}</div>
                )}
              </div>
            </div>
          )}

          {isOwnProfile && pendingRequests.length > 0 && (
            <div className="card" style={{ border: "1px solid #0866ff", background: "#f0f7ff", marginBottom: "24px" }}>
              <h2 style={{ fontSize: "16px", fontWeight: "600", marginBottom: "12px" }}>Follow Requests</h2>
              {pendingRequests.map((req) => (
                <div key={req.id} style={{ display: "flex", justifyContent: "space-between", alignItems: "center", padding: "8px 0" }}>
                <Link href={`/profile/${req.follower_id}`} style={{ display: "flex", alignItems: "center", gap: "8px", textDecoration: "none", color: "inherit", cursor: "pointer" }}>
                  <div style={{ width: "32px", height: "32px", borderRadius: "50%", background: "#ddd" }} />
                  <span style={{ fontSize: "14px", fontWeight: "600" }}>User {req.follower_id}</span>
                </Link>
                  <div style={{ display: "flex", gap: "8px" }}>
                    <button onClick={() => handleRequest(req.id, "accept")} className="btn-primary" style={{ padding: "4px 12px", fontSize: "13px" }}>Confirm</button>
                    <button onClick={() => handleRequest(req.id, "decline")} className="btn-secondary" style={{ padding: "4px 12px", fontSize: "13px" }}>Delete</button>
                  </div>
                </div>
              ))}
            </div>
          )}

          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "16px", marginBottom: "24px" }}>
            <div className="card">
              <h3 style={{ fontSize: "16px", fontWeight: "600", marginBottom: "12px" }}>Followers</h3>
              {followers.length === 0
                ? <p style={{ fontSize: "14px", color: "var(--text-secondary)" }}>No followers yet.</p>
                : followers.map((f) => (
                  <Link key={f.id} href={`/profile/${f.follower_id}`} style={{ display: "flex", alignItems: "center", gap: "8px", padding: "8px 0", textDecoration: "none", color: "inherit", cursor: "pointer" }}>
                    <div style={{ width: "32px", height: "32px", borderRadius: "50%", background: "#ddd" }} />
                    <span style={{ fontSize: "14px" }}>User {f.follower_id}</span>
                  </Link>
                ))
              }
            </div>

            <div className="card">
              <h3 style={{ fontSize: "16px", fontWeight: "600", marginBottom: "12px" }}>Following</h3>
              {following.length === 0
                ? <p style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Not following anyone yet.</p>
                : following.map((f) => (
                  <Link key={f.id} href={`/profile/${f.following_id}`} style={{ display: "flex", alignItems: "center", gap: "8px", padding: "8px 0", textDecoration: "none", color: "inherit", cursor: "pointer" }}>
                    <div style={{ width: "32px", height: "32px", borderRadius: "50%", background: "#ddd" }} />
                    <span style={{ fontSize: "14px" }}>User {f.following_id}</span>
                  </Link>
                ))
              }
            </div>
          </div>

          <div>
            <h2 style={{ fontSize: "18px", fontWeight: "700", marginBottom: "16px" }}>Posts</h2>
            {posts.length === 0
              ? <div className="card" style={{ textAlign: "center", padding: "32px" }}><p style={{ color: "var(--text-secondary)" }}>No posts yet.</p></div>
              : posts.map((post) => <PostCard key={post.id} post={post} />)
            }
          </div>
        </>
      )}
    </div>
  );
}
