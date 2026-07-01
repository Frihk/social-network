'use client';

import { useEffect, useState, useContext } from 'react';
import Link from 'next/link';
import { AuthContext } from '../../context/AuthContext';
import ChatWindow from '../../components/ChatWindow';
import { getDMEligibleUsers } from '../../lib/chat';
import { fetchGroups } from '../../lib/groups';
import NotificationBell from '../../components/NotificationBell';

export default function ChatPage() {
  const { user } = useContext(AuthContext);
  const [selectedConversation, setSelectedConversation] = useState(null);
  const [dmUsers, setDmUsers] = useState([]);
  const [groups, setGroups] = useState([]);
  const [isLoadingDMs, setIsLoadingDMs] = useState(true);
  const [isLoadingGroups, setIsLoadingGroups] = useState(true);
  const [activeTab, setActiveTab] = useState('dms');

  // Fetch DM-eligible users and groups on mount
  useEffect(() => {
    const fetchData = async () => {
      // Fetch DM users
      setIsLoadingDMs(true);
      try {
        const users = await getDMEligibleUsers();
        setDmUsers(users || []);
      } catch (err) {
        console.error('Failed to fetch DM users:', err);
      } finally {
        setIsLoadingDMs(false);
      }

      // Fetch groups
      setIsLoadingGroups(true);
      try {
        const userGroups = await fetchGroups(user?.id);
        setGroups(userGroups || []);
      } catch (err) {
        console.error('Failed to fetch groups:', err);
      } finally {
        setIsLoadingGroups(false);
      }
    };

    fetchData();
  }, []);

  // Handle real-time message updates
  useEffect(() => {
    const handleRealtimeMessage = (event) => {
      const msg = event.detail;
      if (!msg) return;

      // Update DM users list with last message preview
      if (msg.type === 'private_message') {
        setDmUsers(prev => {
          const otherUserId = msg.sender_id === user?.id ? msg.receiver_id : msg.sender_id;
          const index = prev.findIndex(u => u.id === otherUserId);
          if (index >= 0) {
            const newUsers = [...prev];
            // Move to top
            const [dmUser] = newUsers.splice(index, 1);
            newUsers.unshift({
              ...dmUser,
              last_message: msg.content,
              last_message_time: msg.created_at,
            });
            return newUsers;
          }
          return prev;
        });
      }

      // Update groups list with last message preview
      if (msg.type === 'group_message') {
        setGroups(prev => {
          const index = prev.findIndex(g => g.id === msg.group_id);
          if (index >= 0) {
            const newGroups = [...prev];
            // Move to top
            const [group] = newGroups.splice(index, 1);
            newGroups.unshift({
              ...group,
              last_message: msg.content,
              last_message_time: msg.created_at,
            });
            return newGroups;
          }
          return prev;
        });
      }
    };

    window.addEventListener('realtime-message', handleRealtimeMessage);
    return () => window.removeEventListener('realtime-message', handleRealtimeMessage);
  }, [user?.id]);

  const handleSelectConversation = (type, id, name) => {
    setSelectedConversation({ type, id, name });
  };

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Left sidebar - Conversation list */}
      <div className="w-1/4 bg-white border-r border-gray-200 flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b">
          <div className="flex items-center gap-3">
            <Link href="/" className="text-blue-500 hover:text-blue-700 font-bold text-xl">&larr;</Link>
            <h1 className="text-xl font-bold">Messages</h1>
          </div>
          <div className="flex items-center gap-2">
            <NotificationBell />
            <button
            onClick={() => {
              // Refresh data
              getDMEligibleUsers().then(setDmUsers).catch(console.error);
              fetchGroups(user?.id).then(setGroups).catch(console.error);
            }}
            className="text-gray-500 hover:text-gray-700"
            title="Refresh"
          >
            ↻
          </button>
          </div>
        </div>

        {/* Tabs */}
        <div className="flex border-b">
          <button
            onClick={() => setActiveTab('dms')}
            className={`flex-1 py-2 text-center ${
              activeTab === 'dms'
                ? 'border-b-2 border-blue-500 text-blue-500'
                : 'text-gray-500 hover:text-gray-700'
            }`}
          >
            Direct Messages
          </button>
          <button
            onClick={() => setActiveTab('groups')}
            className={`flex-1 py-2 text-center ${
              activeTab === 'groups'
                ? 'border-b-2 border-blue-500 text-blue-500'
                : 'text-gray-500 hover:text-gray-700'
            }`}
          >
            Groups
          </button>
        </div>

        {/* Conversation list */}
        <div className="flex-1 overflow-y-auto">
          {activeTab === 'dms' ? (
            isLoadingDMs ? (
              <div className="p-4 space-y-3">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="animate-pulse">
                    <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                    <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                  </div>
                ))}
              </div>
            ) : dmUsers.length === 0 ? (
              <div className="flex items-center justify-center h-full text-gray-500">
                No conversations yet
              </div>
            ) : (
              dmUsers.map((u) => (
                <div
                  key={u.id}
                  onClick={() => handleSelectConversation('private', u.id, `${u.first_name} ${u.last_name}`)}
                  className={`p-4 cursor-pointer hover:bg-gray-50 border-b ${
                    selectedConversation?.type === 'private' && selectedConversation?.id === u.id
                      ? 'bg-blue-50 border-l-4 border-l-blue-500'
                      : ''
                  }`}
                >
                  <div className="flex items-center gap-3">
                    {u.avatar_path ? (
                      <img
                        src={u.avatar_path}
                        alt="Avatar"
                        className="w-10 h-10 rounded-full object-cover"
                      />
                    ) : (
                      <div className="w-10 h-10 rounded-full bg-gray-300 flex items-center justify-center">
                        <span className="text-gray-600 font-bold">
                          {u.first_name?.[0]}{u.last_name?.[0]}
                        </span>
                      </div>
                    )}
                    <div className="flex-1 min-w-0">
                      <p className="font-medium truncate">
                        {u.first_name} {u.last_name}
                      </p>
                      {u.nickname && (
                        <p className="text-sm text-gray-500 truncate">@{u.nickname}</p>
                      )}
                      {u.last_message && (
                        <p className="text-sm text-gray-400 truncate mt-1">{u.last_message}</p>
                      )}
                    </div>
                  </div>
                </div>
              ))
            )
          ) : (
            isLoadingGroups ? (
              <div className="p-4 space-y-3">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="animate-pulse">
                    <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                    <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                  </div>
                ))}
              </div>
            ) : groups.length === 0 ? (
              <div className="flex items-center justify-center h-full text-gray-500">
                No groups yet
              </div>
            ) : (
              groups.map((g) => (
                <div
                  key={g.id}
                  onClick={() => handleSelectConversation('group', g.id, g.title || g.name)}
                  className={`p-4 cursor-pointer hover:bg-gray-50 border-b ${
                    selectedConversation?.type === 'group' && selectedConversation?.id === g.id
                      ? 'bg-blue-50 border-l-4 border-l-blue-500'
                      : ''
                  }`}
                >
                  <div className="flex-1 min-w-0">
                    <p className="font-medium truncate">{g.title || g.name}</p>
                    {g.description && (
                      <p className="text-sm text-gray-500 truncate">
                        {g.description.length > 50 ? `${g.description.substring(0, 50)}...` : g.description}
                      </p>
                    )}
                    <p className="text-xs text-gray-400 mt-1">
                      {g.member_count} member{g.member_count !== 1 ? 's' : ''}
                    </p>
                    {g.last_message && (
                      <p className="text-sm text-gray-400 truncate mt-1">{g.last_message}</p>
                    )}
                  </div>
                </div>
              ))
            )
          )}
        </div>
      </div>

      {/* Right content area - ChatWindow */}
      <div className="flex-1 flex flex-col">
        {selectedConversation ? (
          <ChatWindow
            conversationType={selectedConversation.type}
            conversationId={selectedConversation.id}
            conversationName={selectedConversation.name}
            onClose={() => setSelectedConversation(null)}
          />
        ) : (
          <div className="flex items-center justify-center h-full text-gray-500">
            <p className="text-lg">Select a conversation to start messaging</p>
          </div>
        )}
      </div>
    </div>
  );
}
