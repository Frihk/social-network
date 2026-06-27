'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import GroupCard from '../../components/GroupCard';
import { fetchGroups, createGroup, requestToJoinGroup } from '../../lib/groups';

export default function GroupsPage() {
  const [groups, setGroups] = useState([]);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [formData, setFormData] = useState({ title: '', description: '' });
  const [loading, setLoading] = useState(true);
  
  // For development, we'll use a mocked user. In a real app this would come from AuthContext.
  const currentUserId = 'user123'; 

  useEffect(() => {
    loadGroups();
  }, []);

  const loadGroups = async () => {
    setLoading(true);
    try {
      const data = await fetchGroups(currentUserId);
      setGroups(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateGroup = async (e) => {
    e.preventDefault();
    try {
      await createGroup(currentUserId, formData);
      setShowCreateModal(false);
      setFormData({ title: '', description: '' });
      loadGroups(); // Refresh list
    } catch (err) {
      console.error(err);
      alert('Failed to create group');
    }
  };

  const handleJoinRequest = async (groupId) => {
    try {
      await requestToJoinGroup(currentUserId, groupId);
      alert('Join request sent!');
    } catch (err) {
      console.error(err);
      alert('Failed to send join request');
    }
  };

  return (
    <div style={styles.container}>
      <div style={{ marginBottom: '16px' }}>
        <Link href="/" style={{ color: '#3b82f6', textDecoration: 'none', fontWeight: '500' }}>&larr; Back to Feed</Link>
      </div>
      <div style={styles.header}>
        <h1 style={styles.title}>Explore Groups</h1>
        <button style={styles.createButton} onClick={() => setShowCreateModal(true)}>
          + Create Group
        </button>
      </div>

      {loading ? (
        <p style={styles.loading}>Loading groups...</p>
      ) : (
        <div style={styles.grid}>
          {groups.length === 0 ? (
            <p style={styles.empty}>No groups found. Be the first to create one!</p>
          ) : (
            groups.map(group => (
              <GroupCard 
                key={group.id} 
                group={group} 
                onJoinRequest={handleJoinRequest} 
                currentUserId={currentUserId}
              />
            ))
          )}
        </div>
      )}

      {showCreateModal && (
        <div style={styles.modalOverlay}>
          <div style={styles.modalContent}>
            <div style={styles.modalHeader}>
              <h2 style={styles.modalTitle}>Create New Group</h2>
              <button style={styles.closeButton} onClick={() => setShowCreateModal(false)}>✕</button>
            </div>
            <form onSubmit={handleCreateGroup} style={styles.form}>
              <div style={styles.formGroup}>
                <label style={styles.label}>Group Name</label>
                <input 
                  type="text" 
                  required
                  style={styles.input}
                  value={formData.title}
                  onChange={(e) => setFormData({...formData, title: e.target.value})}
                  placeholder="e.g. Next.js Developers"
                />
              </div>
              <div style={styles.formGroup}>
                <label style={styles.label}>Description</label>
                <textarea 
                  required
                  style={{ ...styles.input, minHeight: '100px' }}
                  value={formData.description}
                  onChange={(e) => setFormData({...formData, description: e.target.value})}
                  placeholder="What is this group about?"
                />
              </div>
              <div style={styles.formActions}>
                <button type="button" style={styles.cancelButton} onClick={() => setShowCreateModal(false)}>
                  Cancel
                </button>
                <button type="submit" style={styles.submitButton}>
                  Create
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

const styles = {
  container: {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '24px'
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '32px'
  },
  title: {
    fontSize: '2rem',
    fontWeight: 'bold',
    color: '#111827',
    margin: 0
  },
  createButton: {
    backgroundColor: '#10b981',
    color: '#ffffff',
    border: 'none',
    padding: '10px 20px',
    borderRadius: '8px',
    fontSize: '1rem',
    fontWeight: '600',
    cursor: 'pointer',
    transition: 'background-color 0.2s'
  },
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
    gap: '24px'
  },
  loading: {
    textAlign: 'center',
    color: '#6b7280',
    fontSize: '1.1rem',
    marginTop: '40px'
  },
  empty: {
    textAlign: 'center',
    color: '#6b7280',
    gridColumn: '1 / -1',
    marginTop: '40px',
    fontSize: '1.1rem'
  },
  modalOverlay: {
    position: 'fixed',
    top: 0, left: 0, right: 0, bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    zIndex: 1000
  },
  modalContent: {
    backgroundColor: '#ffffff',
    borderRadius: '12px',
    width: '100%',
    maxWidth: '500px',
    padding: '24px',
    boxShadow: '0 20px 25px -5px rgba(0,0,0,0.1)'
  },
  modalHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '20px'
  },
  modalTitle: {
    margin: 0,
    fontSize: '1.5rem',
    color: '#111827'
  },
  closeButton: {
    background: 'none',
    border: 'none',
    fontSize: '1.5rem',
    cursor: 'pointer',
    color: '#6b7280'
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
    gap: '16px'
  },
  formGroup: {
    display: 'flex',
    flexDirection: 'column',
    gap: '6px'
  },
  label: {
    fontSize: '0.9rem',
    fontWeight: '500',
    color: '#374151'
  },
  input: {
    padding: '12px',
    borderRadius: '8px',
    border: '1px solid #d1d5db',
    fontSize: '1rem',
    fontFamily: 'inherit'
  },
  formActions: {
    display: 'flex',
    justifyContent: 'flex-end',
    gap: '12px',
    marginTop: '16px'
  },
  cancelButton: {
    backgroundColor: '#ffffff',
    border: '1px solid #d1d5db',
    color: '#374151',
    padding: '10px 16px',
    borderRadius: '8px',
    fontWeight: '500',
    cursor: 'pointer'
  },
  submitButton: {
    backgroundColor: '#3b82f6',
    border: 'none',
    color: '#ffffff',
    padding: '10px 20px',
    borderRadius: '8px',
    fontWeight: '500',
    cursor: 'pointer'
  }
};