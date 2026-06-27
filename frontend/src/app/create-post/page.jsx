'use client';

import { useAuth } from '@/context/AuthContext';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import PostForm from '@/components/PostForm';
import NotificationBell from '@/components/NotificationBell';

export default function CreatePostPage() {
  const { user, logout, loading } = useAuth();
  const router = useRouter();

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  if (!user && !loading) {
    router.push('/login');
    return null;
  }

  const handleLogout = async () => {
    await logout();
    router.push('/login');
  };

  const handlePostCreated = () => {
    router.push('/');
  };

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Navbar */}
      <nav className="bg-white shadow-md">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <Link href="/" className="text-2xl font-bold text-blue-600 cursor-pointer">Social Network</Link>
          <div className="flex items-center gap-4">
            {user && (
              <>
                <NotificationBell />
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
            )}
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto">
          <div className="mb-4">
            <Link href="/" className="text-blue-600 hover:underline font-medium">&larr; Back to Feed</Link>
          </div>
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-bold mb-4">Create a New Post</h2>
            <PostForm onPostCreated={handlePostCreated} />
          </div>
        </div>
      </div>
    </div>
  );
}
