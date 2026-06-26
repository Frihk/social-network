<<<<<<< HEAD
//
=======
'use client';

import { useAuth } from '@/context/AuthContext';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function HomePage() {
  const { user, logout, loading } = useAuth();
  const router = useRouter();

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
            {user && (
              <>
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
          <div className="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 className="text-xl font-bold mb-4">Welcome, {user?.first_name}!</h2>
            <p className="text-gray-600">
              This is your social network feed. Start by exploring profiles, creating posts,
              and connecting with friends.
            </p>
          </div>

          {/* Posts Section */}
          <div className="bg-white rounded-lg shadow-md p-6">
            <h3 className="text-lg font-bold mb-4">Recent Posts</h3>
            <p className="text-gray-500">No posts yet. Be the first to post!</p>
          </div>
        </div>
      </div>
    </div>
  );
}
>>>>>>> kagunda
