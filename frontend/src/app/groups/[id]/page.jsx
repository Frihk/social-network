'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import EventWidget from '@/components/EventWidget';
import PostCard from '@/components/PostCard';
import PostForm from '@/components/PostForm';
import { getGroupPosts } from '@/lib/posts';
import { 
  fetchGroupDetails, 
  fetchGroupEvents, 
  inviteUserToGroup, 
  acceptMemberRequest, 
  declineMemberRequest 
} from '@/lib/groups';

export default function GroupDetailPage({ params }) {
  const { id } = params;
  const [group, setGroup] = useState(null);
  const [events, setEvents] = useState([]);
  const [posts, setPosts] = useState([]);
  const [postsLoading, setPostsLoading] = useState(true);
  const [inviteUserId, setInviteUserId] = useState('');
  const [loading, setLoading] = useState(true);
  const { user, loading: authLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login');
    } else if (user && id) {
      loadData();
      loadPosts();
    }
  }, [id, authLoading, user]);

  const loadData = async () => {
    setLoading(true);
    try {
      const [groupData, eventsData] = await Promise.all([
        fetchGroupDetails(id).catch(() => null),
        fetchGroupEvents(id).catch(() => [])
      ]);
      setGroup(groupData);
      setEvents(eventsData || []);
    } catch (err) {
      console.error('Error loading group data:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadPosts = async () => {
    setPostsLoading(true);
    try {
      const data = await getGroupPosts(id);
      setPosts(data || []);
    } catch (err) {
      console.error('Error loading group posts:', err);
    } finally {
      setPostsLoading(false);
    }
  };

  const handleInvite = async (e) => {
    e.preventDefault();
    try {
      await inviteUserToGroup(id, inviteUserId);
      setInviteUserId('');
      alert('User invited successfully!');
    } catch (err) {
      console.error(err);
      alert('Failed to invite user');
    }
  };

  const handleAcceptRequest = async (memberId) => {
    try {
      await acceptMemberRequest(id, memberId);
      loadData();
    } catch (err) {
      console.error(err);
      alert('Failed to accept request');
    }
  };

  const handleDeclineRequest = async (memberId) => {
    try {
      await declineMemberRequest(id, memberId);
      loadData();
    } catch (err) {
      console.error(err);
      alert('Failed to decline request');
    }
  };

  if (authLoading || loading || (!user && !authLoading)) {
    return <div className="min-h-screen bg-gray-100 flex items-center justify-center text-gray-500 text-lg">Loading group details...</div>;
  }

  if (!group) {
    return <div className="min-h-screen bg-gray-100 flex items-center justify-center text-red-500 text-lg">Group not found or you do not have access.</div>;
  }

  const isCreator = group.creator_id === user?.id;
  const members = group.members || [];
  const acceptedMembers = members.filter(m => m.status === 'accepted' || m.user_id === group.creator_id);
  const pendingRequests = members.filter(m => m.status === 'requested');

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="max-w-6xl mx-auto px-4 py-6">
        <Link href="/groups" className="text-blue-600 font-semibold hover:underline inline-block mb-4">
          &larr; Back to Groups
        </Link>

        {/* Header Section */}
        <div className="bg-white rounded-lg shadow-md p-8 mb-6">
          <h1 className="text-3xl font-bold text-gray-900 mb-3">{group.title}</h1>
          <p className="text-gray-700 leading-relaxed mb-3">{group.description}</p>
          <p className="text-sm text-gray-500 mb-4">
            Created by: {group.creator_id} &bull; {acceptedMembers.length} Members
          </p>

          {isCreator && (
            <form onSubmit={handleInvite} className="flex gap-3 pt-5 border-t border-gray-100">
              <input
                type="text"
                placeholder="Enter User ID to invite..."
                value={inviteUserId}
                onChange={(e) => setInviteUserId(e.target.value)}
                className="flex-1 px-4 py-2 rounded-lg border border-gray-300 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              />
              <button type="submit" className="bg-blue-600 text-white px-5 py-2 rounded-lg font-semibold hover:bg-blue-700 transition-colors">
                Invite
              </button>
            </form>
          )}
        </div>

        <div className="flex flex-col lg:flex-row gap-6">
          {/* Left Column: Feed */}
          <div className="flex-[2] min-w-0">
            <section>
              <h2 className="text-xl font-bold text-gray-900 mb-4">Post Feed</h2>

              <div className="mb-6">
                <PostForm onPostCreated={loadPosts} groupId={id} />
              </div>

              {postsLoading ? (
                <div className="bg-white rounded-lg shadow-md p-6 text-gray-500">Loading posts...</div>
              ) : posts.length === 0 ? (
                <div className="bg-white rounded-lg shadow-md p-8 text-center text-gray-500">
                  No posts in this group yet. Be the first to post!
                </div>
              ) : (
                <div className="space-y-4">
                  {posts.map(post => (
                    <PostCard key={post.id} post={post} />
                  ))}
                </div>
              )}
            </section>
          </div>

          {/* Right Column: Widgets */}
          <div className="flex-1 min-w-0 flex flex-col gap-6">
            {/* Pending Requests (Creator Only) */}
            {isCreator && pendingRequests.length > 0 && (
              <div className="bg-white rounded-lg shadow-md p-5">
                <h3 className="text-lg font-semibold text-gray-900 mb-3">Pending Requests</h3>
                <ul className="space-y-2">
                  {pendingRequests.map(req => (
                    <li key={req.user_id} className="flex justify-between items-center py-1">
                      <span className="text-sm text-gray-700">{req.user_id}</span>
                      <div className="flex gap-2">
                        <button onClick={() => handleAcceptRequest(req.user_id)} className="bg-green-500 text-white w-7 h-7 rounded flex items-center justify-center text-sm hover:bg-green-600 transition-colors">✓</button>
                        <button onClick={() => handleDeclineRequest(req.user_id)} className="bg-red-500 text-white w-7 h-7 rounded flex items-center justify-center text-sm hover:bg-red-600 transition-colors">✕</button>
                      </div>
                    </li>
                  ))}
                </ul>
              </div>
            )}

            {/* Members List */}
            <div className="bg-white rounded-lg shadow-md p-5">
              <h3 className="text-lg font-semibold text-gray-900 mb-3">Members</h3>
              <ul className="space-y-2">
                {acceptedMembers.slice(0, 10).map(member => (
                  <li key={member.user_id} className="flex items-center gap-2 text-sm">
                    <span>👤 {member.user_id}</span>
                    {member.user_id === group.creator_id && (
                      <span className="bg-gray-100 text-gray-600 px-2 py-0.5 rounded-full text-xs font-medium">Creator</span>
                    )}
                  </li>
                ))}
                {acceptedMembers.length > 10 && (
                  <p className="text-sm text-gray-500 italic text-center pt-2">And {acceptedMembers.length - 10} more...</p>
                )}
              </ul>
            </div>

            {/* Events Widget */}
            <EventWidget
              events={events}
              groupId={id}
              currentUserId={user?.id}
              onEventUpdated={loadData}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
