'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { getUserProfile, updateProfilePrivacy } from '@/lib/auth';

export default function ProfilePage() {
  const { id } = useParams();
  const { user } = useAuth();
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [updatingPrivacy, setUpdatingPrivacy] = useState(false);

  const isOwnProfile = user && user.id === id;

  useEffect(() => {
    loadProfile();
  }, [id]);

  const loadProfile = async () => {
    try {
      const data = await getUserProfile(id);
      setProfile(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handlePrivacyToggle = async () => {
    if (!isOwnProfile) return;

    setUpdatingPrivacy(true);
    try {
      await updateProfilePrivacy(id, !profile.is_private);
      setProfile(prev => ({ ...prev, is_private: !prev.is_private }));
    } catch (err) {
      setError(err.message);
    } finally {
      setUpdatingPrivacy(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">Loading profile...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl text-red-500">{error}</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100 py-8">
      <div className="container mx-auto px-4 max-w-4xl">
        <div className="bg-white rounded-lg shadow-md p-6">
          {/* Profile Header */}
          <div className="flex items-start gap-6 mb-6">
            <div className="w-24 h-24 rounded-full bg-gray-300 flex items-center justify-center overflow-hidden">
              {profile.avatar ? (
                <img src={profile.avatar} alt="Avatar" className="w-full h-full object-cover" />
              ) : (
                <span className="text-3xl text-gray-600">
                  {profile.first_name?.[0]}{profile.last_name?.[0]}
                </span>
              )}
            </div>

            <div className="flex-1">
              <h1 className="text-3xl font-bold mb-2">
                {profile.first_name} {profile.last_name}
              </h1>
              {profile.nickname && (
                <p className="text-gray-600 mb-2">@{profile.nickname}</p>
              )}
              {profile.about_me && (
                <p className="text-gray-700 mb-4">{profile.about_me}</p>
              )}

              {isOwnProfile && (
                <div className="flex items-center gap-4">
                  <button
                    onClick={handlePrivacyToggle}
                    disabled={updatingPrivacy}
                    className={`px-4 py-2 rounded-lg ${
                      profile.is_private
                        ? 'bg-yellow-500 hover:bg-yellow-600'
                        : 'bg-green-500 hover:bg-green-600'
                    } text-white disabled:bg-gray-400`}
                  >
                    {updatingPrivacy
                      ? 'Updating...'
                      : profile.is_private
                      ? '🔒 Private Profile'
                      : '🌍 Public Profile'}
                  </button>
                </div>
              )}

              {!isOwnProfile && profile.is_private && (
                <div className="text-yellow-600">
                  🔒 This is a private profile
                </div>
              )}
            </div>
          </div>

          {/* User Information */}
          {profile.email && (
            <div className="border-t pt-6">
              <h2 className="text-xl font-bold mb-4">Information</h2>
              <div className="space-y-2">
                <div>
                  <span className="font-semibold">Email:</span> {profile.email}
                </div>
                {profile.date_of_birth && (
                  <div>
                    <span className="font-semibold">Date of Birth:</span> {profile.date_of_birth}
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Posts Section */}
          <div className="border-t pt-6 mt-6">
            <h2 className="text-xl font-bold mb-4">Posts</h2>
            <p className="text-gray-500">No posts yet</p>
          </div>

          {/* Followers Section */}
          <div className="border-t pt-6 mt-6">
            <h2 className="text-xl font-bold mb-4">Followers & Following</h2>
            <div className="flex gap-8">
              <div>
                <span className="font-semibold">Followers:</span> 0
              </div>
              <div>
                <span className="font-semibold">Following:</span> 0
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
