'use client';

import React from 'react';
import Link from 'next/link';

export default function GroupCard({ group, onJoinRequest, currentUserId }) {
  // If the user is the creator or already an accepted member
  // Note: The backend group details might include 'members' if the user is accepted.
  // In the general groups list, the backend might just return basic info. We'll handle state gracefully.
  const isCreator = group.creator_id === currentUserId;
  const isMember = isCreator || (group.members && group.members.some(m => m.user_id === currentUserId && m.status === 'accepted'));
  
  // Basic truncation for description
  const truncateDesc = (text, length = 100) => {
    if (!text) return '';
    return text.length > length ? text.substring(0, length) + '...' : text;
  };

  const memberCount = group.members ? group.members.length : (group.member_count || 1);

  return (
    <div style={styles.card}>
      <div style={styles.content}>
        <h3 style={styles.title}>
          <Link href={`/groups/${group.id}`} style={styles.link}>
            {group.title}
          </Link>
        </h3>
        <p style={styles.description}>{truncateDesc(group.description)}</p>
        <div style={styles.footer}>
          <span style={styles.memberCount}>
            👥 {memberCount} {memberCount === 1 ? 'Member' : 'Members'}
          </span>
          {isMember ? (
            <button style={{ ...styles.button, ...styles.memberButton }} disabled>
              Member
            </button>
          ) : (
            <button 
              style={{ ...styles.button, ...styles.joinButton }} 
              onClick={() => onJoinRequest(group.id)}
            >
              Request to Join
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

const styles = {
  card: {
    border: '1px solid #eaeaea',
    borderRadius: '12px',
    padding: '20px',
    backgroundColor: '#ffffff',
    boxShadow: '0 2px 8px rgba(0,0,0,0.05)',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-between',
    transition: 'transform 0.2s, box-shadow 0.2s',
    marginBottom: '16px'
  },
  content: {
    display: 'flex',
    flexDirection: 'column',
    gap: '12px'
  },
  title: {
    margin: 0,
    fontSize: '1.25rem',
    fontWeight: '600',
    color: '#111827'
  },
  link: {
    textDecoration: 'none',
    color: 'inherit'
  },
  description: {
    margin: 0,
    color: '#4B5563',
    fontSize: '0.95rem',
    lineHeight: '1.5'
  },
  footer: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: '12px',
    paddingTop: '12px',
    borderTop: '1px solid #f3f4f6'
  },
  memberCount: {
    fontSize: '0.875rem',
    color: '#6B7280',
    fontWeight: '500'
  },
  button: {
    padding: '8px 16px',
    borderRadius: '6px',
    fontSize: '0.875rem',
    fontWeight: '500',
    cursor: 'pointer',
    border: 'none',
    transition: 'background-color 0.2s'
  },
  joinButton: {
    backgroundColor: '#3b82f6',
    color: '#ffffff'
  },
  memberButton: {
    backgroundColor: '#f3f4f6',
    color: '#9ca3af',
    cursor: 'not-allowed'
  }
};