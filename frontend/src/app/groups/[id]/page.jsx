'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import EventWidget from '../../../components/EventWidget';
import { 
  fetchGroupDetails, 
  fetchGroupEvents, 
  inviteUserToGroup, 
  acceptMemberRequest, 
  declineMemberRequest 
} from '../../../lib/groups';

export default function GroupDetailPage({ params }) {
  const { id } = params;
  const [group, setGroup] = useState(null);
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [inviteUserId, setInviteUserId] = useState('');
  
  // For development, we mock the current user.
  const currentUserId = 'user123';

  useEffect(() => {
    if (id) {
      loadData();
    }
  }, [id]);

  const loadData = async () => {
    setLoading(true);
    try {
      const [groupData, eventsData] = await Promise.all([
        fetchGroupDetails(currentUserId, id).catch(() => null),
        fetchGroupEvents(currentUserId, id).catch(() => [])
      ]);
      setGroup(groupData);
      setEvents(eventsData || []);
    } catch (err) {
      console.error('Error loading group data:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleInvite = async (e) => {
    e.preventDefault();
    try {
      await inviteUserToGroup(currentUserId, id, inviteUserId);
      setInviteUserId('');
      alert('User invited successfully!');
    } catch (err) {
      console.error(err);
      alert('Failed to invite user');
    }
  };

  const handleAcceptRequest = async (memberId) => {
    try {
      await acceptMemberRequest(currentUserId, id, memberId);
      loadData(); // Refresh group details
    } catch (err) {
      console.error(err);
      alert('Failed to accept request');
    }
  };

  const handleDeclineRequest = async (memberId) => {
    try {
      await declineMemberRequest(currentUserId, id, memberId);
      loadData(); // Refresh group details
    } catch (err) {
      console.error(err);
      alert('Failed to decline request');
    }
  };

  if (loading) {
    return <div style={styles.loading}>Loading group details...</div>;
  }

  if (!group) {
    return <div style={styles.error}>Group not found or you do not have access.</div>;
  }

  const isCreator = group.creator_id === currentUserId;
  // Fallback to empty array if members is undefined
  const members = group.members || [];
  const acceptedMembers = members.filter(m => m.status === 'accepted' || m.user_id === group.creator_id);
  const pendingRequests = members.filter(m => m.status === 'requested');

  return (
    <div style={styles.container}>
      <div style={{ marginBottom: '16px' }}>
        <Link href="/groups" style={{ color: '#3b82f6', textDecoration: 'none', fontWeight: '500' }}>&larr; Back to Groups</Link>
      </div>
      {/* Header Section */}
      <div style={styles.headerCard}>
        <h1 style={styles.title}>{group.title}</h1>
        <p style={styles.description}>{group.description}</p>
        <p style={styles.meta}>
          Created by: {group.creator_id} • {acceptedMembers.length} Members
        </p>
        
        {isCreator && (
          <form onSubmit={handleInvite} style={styles.inviteForm}>
            <input 
              type="text" 
              placeholder="Enter User ID to invite..." 
              value={inviteUserId}
              onChange={(e) => setInviteUserId(e.target.value)}
              style={styles.input}
              required
            />
            <button type="submit" style={styles.inviteButton}>Invite</button>
          </form>
        )}
      </div>

      <div style={styles.layout}>
        {/* Left Column: Feed */}
        <div style={styles.mainCol}>
          <div style={styles.section}>
            <h2 style={styles.sectionTitle}>Post Feed</h2>
            <div style={styles.placeholderCard}>
              <p>Write something...</p>
              <button style={styles.mockButton}>Post</button>
            </div>
            {/* Mock Posts */}
            <div style={{...styles.placeholderCard, marginTop: '16px'}}>
              <strong>User456</strong>
              <p>Hello everyone, excited to be here!</p>
            </div>
            <div style={{...styles.placeholderCard, marginTop: '16px'}}>
              <strong>User789</strong>
              <p>Has anyone seen the latest updates?</p>
            </div>
          </div>
        </div>

        {/* Right Column: Widgets */}
        <div style={styles.sideCol}>
          
          {/* Pending Requests (Creator Only) */}
          {isCreator && pendingRequests.length > 0 && (
            <div style={styles.widget}>
              <h3 style={styles.widgetTitle}>Pending Requests</h3>
              <ul style={styles.list}>
                {pendingRequests.map(req => (
                  <li key={req.user_id} style={styles.listItem}>
                    <span>{req.user_id}</span>
                    <div style={styles.actionButtons}>
                      <button onClick={() => handleAcceptRequest(req.user_id)} style={styles.acceptBtn}>✓</button>
                      <button onClick={() => handleDeclineRequest(req.user_id)} style={styles.declineBtn}>✕</button>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          )}

          {/* Members List */}
          <div style={styles.widget}>
            <h3 style={styles.widgetTitle}>Members</h3>
            <ul style={styles.list}>
              {acceptedMembers.slice(0, 10).map(member => (
                <li key={member.user_id} style={styles.listItem}>
                  <span>👤 {member.user_id}</span>
                  {member.user_id === group.creator_id && <span style={styles.badge}>Creator</span>}
                </li>
              ))}
              {acceptedMembers.length > 10 && (
                <p style={styles.moreText}>And {acceptedMembers.length - 10} more...</p>
              )}
            </ul>
          </div>

          {/* Events Widget */}
          <EventWidget 
            events={events} 
            groupId={id} 
            currentUserId={currentUserId} 
            onEventUpdated={loadData}
          />

        </div>
      </div>
    </div>
  );
}

const styles = {
  container: {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '24px'
  },
  loading: {
    textAlign: 'center',
    padding: '40px',
    fontSize: '1.2rem',
    color: '#6b7280'
  },
  error: {
    textAlign: 'center',
    padding: '40px',
    fontSize: '1.2rem',
    color: '#ef4444'
  },
  headerCard: {
    backgroundColor: '#ffffff',
    borderRadius: '12px',
    padding: '32px',
    border: '1px solid #eaeaea',
    boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.05)',
    marginBottom: '24px'
  },
  title: {
    fontSize: '2.5rem',
    fontWeight: 'bold',
    color: '#111827',
    margin: '0 0 16px 0'
  },
  description: {
    fontSize: '1.1rem',
    color: '#4b5563',
    lineHeight: '1.6',
    margin: '0 0 16px 0'
  },
  meta: {
    fontSize: '0.9rem',
    color: '#6b7280',
    margin: '0 0 24px 0'
  },
  inviteForm: {
    display: 'flex',
    gap: '12px',
    borderTop: '1px solid #f3f4f6',
    paddingTop: '20px'
  },
  input: {
    padding: '10px 16px',
    borderRadius: '8px',
    border: '1px solid #d1d5db',
    fontSize: '0.95rem',
    minWidth: '250px'
  },
  inviteButton: {
    backgroundColor: '#3b82f6',
    color: '#ffffff',
    border: 'none',
    padding: '10px 20px',
    borderRadius: '8px',
    fontWeight: '500',
    cursor: 'pointer'
  },
  layout: {
    display: 'flex',
    gap: '24px',
    flexDirection: 'row',
    flexWrap: 'wrap'
  },
  mainCol: {
    flex: '2',
    minWidth: '300px'
  },
  sideCol: {
    flex: '1',
    minWidth: '300px',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px'
  },
  section: {
    backgroundColor: 'transparent'
  },
  sectionTitle: {
    fontSize: '1.5rem',
    fontWeight: 'bold',
    color: '#111827',
    marginBottom: '16px'
  },
  placeholderCard: {
    backgroundColor: '#ffffff',
    border: '1px solid #eaeaea',
    borderRadius: '12px',
    padding: '20px',
    boxShadow: '0 1px 3px rgba(0,0,0,0.05)'
  },
  mockButton: {
    marginTop: '12px',
    backgroundColor: '#10b981',
    color: '#ffffff',
    border: 'none',
    padding: '8px 16px',
    borderRadius: '6px',
    cursor: 'pointer'
  },
  widget: {
    backgroundColor: '#ffffff',
    border: '1px solid #eaeaea',
    borderRadius: '12px',
    padding: '20px',
    boxShadow: '0 1px 3px rgba(0,0,0,0.05)'
  },
  widgetTitle: {
    fontSize: '1.2rem',
    fontWeight: '600',
    color: '#111827',
    margin: '0 0 16px 0'
  },
  list: {
    listStyle: 'none',
    padding: 0,
    margin: 0,
    display: 'flex',
    flexDirection: 'column',
    gap: '12px'
  },
  listItem: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingBottom: '8px',
    borderBottom: '1px solid #f3f4f6',
    fontSize: '0.95rem'
  },
  badge: {
    backgroundColor: '#f3f4f6',
    color: '#4b5563',
    padding: '2px 8px',
    borderRadius: '12px',
    fontSize: '0.75rem',
    fontWeight: '500'
  },
  moreText: {
    fontSize: '0.85rem',
    color: '#6b7280',
    fontStyle: 'italic',
    textAlign: 'center'
  },
  actionButtons: {
    display: 'flex',
    gap: '8px'
  },
  acceptBtn: {
    backgroundColor: '#10b981',
    color: '#ffffff',
    border: 'none',
    borderRadius: '4px',
    width: '28px',
    height: '28px',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center'
  },
  declineBtn: {
    backgroundColor: '#ef4444',
    color: '#ffffff',
    border: 'none',
    borderRadius: '4px',
    width: '28px',
    height: '28px',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center'
  }
};