'use client';

import { useCallback, useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import NotificationBell from '@/components/NotificationBell';
import PostCard from '@/components/PostCard';
import { getFeed } from '@/lib/posts';

export default function HomePage() {
  const { user, logout, loading } = useAuth();
  const router = useRouter();
  const [posts, setPosts] = useState([]);
  const [feedLoading, setFeedLoading] = useState(true);
  const [feedError, setFeedError] = useState('');

  const fetchFeed = useCallback(async () => {
    setFeedLoading(true);
    setFeedError('');
    try {
      const data = await getFeed();
      setPosts(data || []);
    } catch (err) {
      setFeedError(err.message);
    } finally {
      setFeedLoading(false);
    }
  }, []);

  useEffect(() => {
    if (!loading) fetchFeed();
  }, [loading, user, fetchFeed]);

  const handleLogout = async () => {
    await logout();
    router.push('/login');
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Navbar */}
      <nav className="bg-white shadow-md">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-blue-600">Social Network</h1>
          <div className="flex items-center gap-4">
            {user ? (
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
                  href={`/profile/${user.id}`}
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

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto">
          {user && (
            <div className="bg-white rounded-lg shadow-md p-6 mb-6 flex justify-between items-center">
              <h2 className="text-xl font-bold">Welcome, {user.first_name}!</h2>
            </div>
          )}

          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-bold">Recent Posts</h3>
              <button
                type="button"
                onClick={fetchFeed}
                className="rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50"
              >
                Refresh
              </button>
            </div>

            {feedError && (
              <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-700">
                {feedError}
              </div>
            )}

            {feedLoading ? (
              <div className="bg-white rounded-lg shadow-md p-6 text-gray-500">Loading feed...</div>
            ) : posts.length === 0 ? (
              <div className="bg-white rounded-lg shadow-md p-6 text-gray-500">No posts yet. Be the first to post!</div>
            ) : (
              posts.map((post) => <PostCard key={post.id} post={post} />)
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
