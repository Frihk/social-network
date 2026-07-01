"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/context/AuthContext";
import { getUserProfile, updateProfilePrivacy, updateUserProfile } from "@/lib/auth";
import { getUserPosts } from "@/lib/posts";
import { getFollowers, getFollowing, followUser, unfollowUser, acceptFollow, declineFollow } from "@/lib/followers";
import PostCard from "@/components/PostCard";
import NotificationBell from "@/components/NotificationBell";

export default function ProfilePage() {
  const { id: profileId } = useParams();
  const { user: currentUser, logout } = useAuth();
  const router = useRouter();
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

  const handleLogout = async () => {
    await logout();
    router.push("/login");
  };

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
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-gray-500 text-lg">Loading profile...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="rounded-lg border border-red-200 bg-red-50 p-6 text-red-700 text-center">
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Navbar */}
      <nav className="bg-white shadow-md">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <Link href="/" className="text-2xl font-bold text-blue-600">Social Network</Link>
          <div className="flex items-center gap-4">
            {currentUser ? (
              <>
                <NotificationBell />
                <Link
                  href="/create-post"
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors font-semibold"
                >
                  + Create Post
                </Link>
                <Link
                  href="/chat"
                  className="text-gray-700 hover:text-blue-600"
                >
                  Chat
                </Link>
                <Link
                  href="/groups"
                  className="text-gray-700 hover:text-blue-600"
                >
                  Groups
                </Link>
                <Link
                  href={`/profile/${currentUser.id}`}
                  className="text-gray-700 hover:text-blue-600"
                >
                  My Profile
                </Link>
                <button
                  onClick={handleLogout}
                  className="bg-red-500 text-white px-4 py-2 rounded-lg hover:bg-red-600"
                >
                  Logout
                </button>
              </>
            ) : (
              <div className="flex items-center gap-2">
                <Link
                  href="/login"
                  className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors font-semibold"
                >
                  Login
                </Link>
                <Link
                  href="/register"
                  className="bg-emerald-600 text-white px-4 py-2 rounded-lg hover:bg-emerald-700 transition-colors font-semibold"
                >
                  Sign Up
                </Link>
              </div>
            )}
          </div>
        </div>
      </nav>

      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto">
          <Link href="/" className="text-blue-600 font-semibold hover:underline inline-block mb-6">
            &larr; Back to Feed
          </Link>

          {/* Profile Header Card */}
          <div className="bg-white rounded-lg shadow-md p-8 text-center mb-6">
            <div className="w-[120px] h-[120px] rounded-full bg-gray-200 mx-auto mb-4 overflow-hidden flex items-center justify-center">
              {profileUser?.avatar ? (
                <img src={profileUser.avatar} alt="avatar" className="w-full h-full object-cover" />
              ) : (
                <span className="text-3xl text-gray-500 font-semibold">
                  {profileUser?.first_name?.[0]}{profileUser?.last_name?.[0]}
                </span>
              )}
            </div>

            <h1 className="text-2xl font-bold text-gray-900">
              {profileUser ? `${profileUser.first_name} ${profileUser.last_name}` : "Unknown User"}
            </h1>

            {profileUser?.nickname && (
              <p className="text-gray-500 mt-1">@{profileUser.nickname}</p>
            )}

            {profileUser?.about_me && (
              <p className="text-gray-700 text-sm mt-2 max-w-md mx-auto">{profileUser.about_me}</p>
            )}

            <div className="flex justify-center gap-6 mt-4 mb-4">
              <div>
                <span className="font-bold text-gray-900">{posts.length}</span>
                <span className="text-gray-500 text-sm ml-1">Posts</span>
              </div>
              <div>
                <span className="font-bold text-gray-900">{followers.length}</span>
                <span className="text-gray-500 text-sm ml-1">Followers</span>
              </div>
              <div>
                <span className="font-bold text-gray-900">{following.length}</span>
                <span className="text-gray-500 text-sm ml-1">Following</span>
              </div>
            </div>

            {isOwnProfile && (
              <div className="mt-4 flex justify-center gap-3">
                <label className="rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50 cursor-pointer inline-block transition-colors">
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
                  className={
                    profileUser?.is_private
                      ? "rounded-lg border border-gray-200 bg-white px-4 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50 transition-colors"
                      : "bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors font-semibold"
                  }
                >
                  {updatingPrivacy
                    ? "Updating..."
                    : profileUser?.is_private
                    ? "Private Profile"
                    : "Public Profile"}
                </button>
              </div>
            )}

            {!isOwnProfile && currentUser && (
              <div className="mt-4">
                <button
                  onClick={handleFollowToggle}
                  disabled={followState === "pending"}
                  className={
                    followState === "following" || followState === "pending"
                      ? "rounded-lg border border-gray-200 bg-white px-6 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50 transition-colors"
                      : "bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition-colors font-semibold"
                  }
                  style={{
                    opacity: followState === "pending" ? 0.6 : 1,
                    cursor: followState === "pending" ? "not-allowed" : "pointer",
                  }}
                >
                  {followState === "following" ? "Following" : followState === "pending" ? "Requested" : "Follow"}
                </button>
              </div>
            )}
          </div>

          {!isOwnProfile && profileUser?.is_private && followState !== "following" ? (
            <div className="bg-white rounded-lg shadow-md p-12 text-center mb-6">
              <p className="text-lg font-semibold text-gray-900 mb-2">This account is private.</p>
              <p className="text-gray-500 text-sm">Follow to see their posts and activity.</p>
            </div>
          ) : (
            <>
              {isOwnProfile && profileUser?.email && (
                <div className="bg-white rounded-lg shadow-md p-6 mb-6">
                  <h2 className="text-base font-semibold text-gray-900 mb-3">Private Information</h2>
                  <div className="text-sm text-gray-700 space-y-1">
                    <div><span className="font-semibold">Email:</span> {profileUser.email}</div>
                    {profileUser.date_of_birth && (
                      <div><span className="font-semibold">Date of Birth:</span> {profileUser.date_of_birth}</div>
                    )}
                  </div>
                </div>
              )}

              {isOwnProfile && pendingRequests.length > 0 && (
                <div className="bg-white rounded-lg shadow-md p-6 mb-6 border-l-4 border-blue-500">
                  <h2 className="text-base font-semibold text-gray-900 mb-3">Follow Requests</h2>
                  <div className="space-y-2">
                    {pendingRequests.map((req) => (
                      <div key={req.id} className="flex justify-between items-center py-2">
                        <Link
                          href={`/profile/${req.follower_id}`}
                          className="flex items-center gap-2 text-gray-900 hover:text-blue-600 no-underline"
                        >
                          <div className="w-8 h-8 rounded-full bg-gray-200 flex-shrink-0" />
                          <span className="text-sm font-semibold">User {req.follower_id}</span>
                        </Link>
                        <div className="flex gap-2">
                          <button
                            onClick={() => handleRequest(req.id, "accept")}
                            className="bg-blue-600 text-white px-3 py-1.5 rounded-lg text-xs font-semibold hover:bg-blue-700 transition-colors"
                          >
                            Confirm
                          </button>
                          <button
                            onClick={() => handleRequest(req.id, "decline")}
                            className="rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-xs font-semibold text-gray-700 hover:bg-gray-50 transition-colors"
                          >
                            Delete
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
                <div className="bg-white rounded-lg shadow-md p-6">
                  <h3 className="text-base font-semibold text-gray-900 mb-3">Followers</h3>
                  {followers.length === 0 ? (
                    <p className="text-sm text-gray-500">No followers yet.</p>
                  ) : (
                    <div className="space-y-2">
                      {followers.map((f) => (
                        <Link
                          key={f.id}
                          href={`/profile/${f.follower_id}`}
                          className="flex items-center gap-2 py-1 text-gray-900 hover:text-blue-600 no-underline"
                        >
                          <div className="w-8 h-8 rounded-full bg-gray-200 flex-shrink-0" />
                          <span className="text-sm">{f.follower_name || `User ${f.follower_id}`}</span>
                        </Link>
                      ))}
                    </div>
                  )}
                </div>

                <div className="bg-white rounded-lg shadow-md p-6">
                  <h3 className="text-base font-semibold text-gray-900 mb-3">Following</h3>
                  {following.length === 0 ? (
                    <p className="text-sm text-gray-500">Not following anyone yet.</p>
                  ) : (
                    <div className="space-y-2">
                      {following.map((f) => (
                        <Link
                          key={f.id}
                          href={`/profile/${f.following_id}`}
                          className="flex items-center gap-2 py-1 text-gray-900 hover:text-blue-600 no-underline"
                        >
                          <div className="w-8 h-8 rounded-full bg-gray-200 flex-shrink-0" />
                          <span className="text-sm">{f.following_name || `User ${f.following_id}`}</span>
                        </Link>
                      ))}
                    </div>
                  )}
                </div>
              </div>

              <div>
                <h2 className="text-lg font-bold text-gray-900 mb-4">Posts</h2>
                {posts.length === 0 ? (
                  <div className="bg-white rounded-lg shadow-md p-8 text-center">
                    <p className="text-gray-500">No posts yet.</p>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {posts.map((post) => (
                      <PostCard key={post.id} post={post} />
                    ))}
                  </div>
                )}
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
