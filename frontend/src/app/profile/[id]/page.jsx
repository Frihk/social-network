"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import {
  getFollowers,
  getFollowing,
  followUser,
  unfollowUser,
  acceptFollow,
  declineFollow,
} from "@/lib/followers";

export default function ProfilePage() {
  const { id: profileId } = useParams();
  const { user: currentUser } = useAuth();
  const isOwnProfile = currentUser && currentUser.id === parseInt(profileId);

  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const [pendingRequests, setPendingRequests] = useState([]);
  const [isFollowing, setIsFollowing] = useState(false);
  const [loading, setLoading] = useState(true);

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const [followersData, followingData] = await Promise.all([
        getFollowers(profileId),
        getFollowing(profileId),
      ]);

      const acceptedFollowers = followersData?.filter((f) => f.status === "accepted") || [];
      const pendingFollowers = followersData?.filter((f) => f.status === "pending") || [];

      setFollowers(acceptedFollowers);
      setFollowing(followingData || []);

      if (isOwnProfile) {
        setPendingRequests(pendingFollowers);
      }

      if (currentUser && !isOwnProfile) {
        const myFollowing = await getFollowing(currentUser.id);
        const followingProfile = myFollowing?.find(
          (f) => f.following_id === parseInt(profileId)
        );
        setIsFollowing(!!followingProfile);
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [profileId, currentUser, isOwnProfile]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleFollowToggle = async () => {
    try {
      if (isFollowing) {
        await unfollowUser(profileId);
      } else {
        await followUser(profileId);
      }
      fetchData();
    } catch (err) {
      console.error(err);
    }
  };

  const handleRequest = async (requestId, action) => {
    try {
      if (action === "accept") {
        await acceptFollow(requestId);
      } else {
        await declineFollow(requestId);
      }
      fetchData();
    } catch (err) {
      console.error(err);
    }
  };

  if (loading) return <div className="container" style={{ textAlign: "center" }}><p>Loading profile...</p></div>;

  return (
    <div className="container">
      <div className="card" style={{ padding: "40px", textAlign: "center", marginBottom: "24px" }}>
        <div style={{ width: "120px", height: "120px", borderRadius: "50%", background: "#ddd", margin: "0 auto 16px" }} />
        <h1 style={{ fontSize: "28px", fontWeight: "700" }}>User {profileId}</h1>
        
        {!isOwnProfile && currentUser && (
          <div style={{ marginTop: "16px" }}>
            <button 
              onClick={handleFollowToggle}
              className={isFollowing ? "btn-secondary" : "btn-primary"}
              style={{ minWidth: "120px" }}
            >
              {isFollowing ? "Following" : "Follow"}
            </button>
          </div>
        )}
      </div>

      {isOwnProfile && pendingRequests.length > 0 && (
        <div className="card" style={{ border: "1px solid #0866ff", background: "#f0f7ff" }}>
          <h2 style={{ fontSize: "16px", fontWeight: "600", marginBottom: "12px" }}>Follow Requests</h2>
          {pendingRequests.map((req) => (
            <div key={req.id} style={{ display: "flex", justifyContent: "space-between", alignItems: "center", padding: "8px 0" }}>
              <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                <div style={{ width: "32px", height: "32px", borderRadius: "50%", background: "#ddd" }} />
                <span style={{ fontSize: "14px", fontWeight: "600" }}>User {req.follower_id}</span>
              </div>
              <div style={{ display: "flex", gap: "8px" }}>
                <button onClick={() => handleRequest(req.id, "accept")} className="btn-primary" style={{ padding: "4px 12px", fontSize: "13px" }}>Confirm</button>
                <button onClick={() => handleRequest(req.id, "decline")} className="btn-secondary" style={{ padding: "4px 12px", fontSize: "13px" }}>Delete</button>
              </div>
            </div>
          ))}
        </div>
      )}

      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "16px" }}>
        <div className="card">
          <h3 style={{ fontSize: "16px", fontWeight: "600", marginBottom: "12px" }}>Followers ({followers.length})</h3>
          {followers.length === 0 ? (
            <p style={{ fontSize: "14px", color: "var(--text-secondary)" }}>No followers yet.</p>
          ) : (
            followers.map((f) => (
              <div key={f.id} style={{ display: "flex", alignItems: "center", gap: "8px", padding: "8px 0" }}>
                <div style={{ width: "32px", height: "32px", borderRadius: "50%", background: "#ddd" }} />
                <span style={{ fontSize: "14px" }}>User {f.follower_id}</span>
              </div>
            ))
          )}
        </div>

        <div className="card">
          <h3 style={{ fontSize: "16px", fontWeight: "600", marginBottom: "12px" }}>Following ({following.length})</h3>
          {following.length === 0 ? (
            <p style={{ fontSize: "14px", color: "var(--text-secondary)" }}>Not following anyone yet.</p>
          ) : (
            following.map((f) => (
              <div key={f.id} style={{ display: "flex", alignItems: "center", gap: "8px", padding: "8px 0" }}>
                <div style={{ width: "32px", height: "32px", borderRadius: "50%", background: "#ddd" }} />
                <span style={{ fontSize: "14px" }}>User {f.following_id}</span>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
