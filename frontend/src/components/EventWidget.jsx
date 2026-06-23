'use client';

import React, { useState } from 'react';
import { createEvent, respondToEvent } from '../lib/groups';

export default function EventWidget({ events, groupId, currentUserId, onEventUpdated }) {
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [formData, setFormData] = useState({ title: '', description: '', event_time: '' });
  const [loading, setLoading] = useState(false);

  const handleCreateEvent = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      await createEvent(currentUserId, groupId, {
        title: formData.title,
        description: formData.description,
        event_time: new Date(formData.event_time).toISOString()
      });
      setShowCreateForm(false);
      setFormData({ title: '', description: '', event_time: '' });
      if (onEventUpdated) onEventUpdated();
    } catch (err) {
      console.error(err);
      alert('Failed to create event');
    } finally {
      setLoading(false);
    }
  };

  const handleRsvp = async (eventId, responseType) => {
    try {
      await respondToEvent(currentUserId, eventId, responseType);
      if (onEventUpdated) onEventUpdated();
    } catch (err) {
      console.error(err);
      alert('Failed to RSVP');
    }
  };

  return (
    <div style={styles.widget}>
      <div style={styles.header}>
        <h3 style={styles.title}>Group Events</h3>
        <button 
          style={styles.createButton}
          onClick={() => setShowCreateForm(!showCreateForm)}
        >
          {showCreateForm ? 'Cancel' : '+ New Event'}
        </button>
      </div>

      {showCreateForm && (
        <form onSubmit={handleCreateEvent} style={styles.form}>
          <input 
            type="text" 
            placeholder="Event Title" 
            required 
            style={styles.input}
            value={formData.title}
            onChange={(e) => setFormData({...formData, title: e.target.value})}
          />
          <textarea 
            placeholder="Description" 
            required 
            style={{ ...styles.input, minHeight: '60px' }}
            value={formData.description}
            onChange={(e) => setFormData({...formData, description: e.target.value})}
          />
          <input 
            type="datetime-local" 
            required 
            style={styles.input}
            value={formData.event_time}
            onChange={(e) => setFormData({...formData, event_time: e.target.value})}
          />
          <button type="submit" disabled={loading} style={styles.submitButton}>
            {loading ? 'Creating...' : 'Create Event'}
          </button>
        </form>
      )}

      <div style={styles.eventList}>
        {(!events || events.length === 0) ? (
          <p style={styles.emptyText}>No upcoming events</p>
        ) : (
          events.map(event => {
            const dateStr = new Date(event.event_time).toLocaleString();
            
            // Note: The backend might return RSVP counts or responses based on implementation.
            // The API guide says "list events ... RSVP counts". Let's assume there is an array `responses`.
            const goingCount = event.responses ? event.responses.filter(r => r.response === 'going').length : 0;
            const notGoingCount = event.responses ? event.responses.filter(r => r.response === 'not_going').length : 0;
            
            // Check user's own response
            const userResponse = event.responses ? event.responses.find(r => r.user_id === currentUserId)?.response : null;

            return (
              <div key={event.id} style={styles.eventItem}>
                <h4 style={styles.eventTitle}>{event.title}</h4>
                <p style={styles.eventDate}>📅 {dateStr}</p>
                <p style={styles.eventDesc}>{event.description}</p>
                
                <div style={styles.rsvpSection}>
                  <div style={styles.rsvpCounts}>
                    <span>✅ {goingCount} Going</span>
                    <span>❌ {notGoingCount} Not Going</span>
                  </div>
                  <div style={styles.rsvpButtons}>
                    <button 
                      style={{
                        ...styles.rsvpButton, 
                        ...(userResponse === 'going' ? styles.activeGoing : {})
                      }}
                      onClick={() => handleRsvp(event.id, 'going')}
                    >
                      Going
                    </button>
                    <button 
                      style={{
                        ...styles.rsvpButton, 
                        ...(userResponse === 'not_going' ? styles.activeNotGoing : {})
                      }}
                      onClick={() => handleRsvp(event.id, 'not_going')}
                    >
                      Not Going
                    </button>
                  </div>
                </div>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}

const styles = {
  widget: {
    backgroundColor: '#ffffff',
    borderRadius: '12px',
    border: '1px solid #eaeaea',
    padding: '20px',
    boxShadow: '0 2px 8px rgba(0,0,0,0.05)',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '16px'
  },
  title: {
    margin: 0,
    fontSize: '1.25rem',
    fontWeight: '600',
    color: '#111827'
  },
  createButton: {
    backgroundColor: '#f3f4f6',
    color: '#374151',
    border: 'none',
    padding: '6px 12px',
    borderRadius: '6px',
    fontSize: '0.875rem',
    fontWeight: '500',
    cursor: 'pointer'
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
    gap: '10px',
    marginBottom: '20px',
    padding: '16px',
    backgroundColor: '#f9fafb',
    borderRadius: '8px',
    border: '1px solid #e5e7eb'
  },
  input: {
    padding: '10px',
    borderRadius: '6px',
    border: '1px solid #d1d5db',
    fontSize: '0.9rem',
    fontFamily: 'inherit'
  },
  submitButton: {
    backgroundColor: '#10b981',
    color: 'white',
    border: 'none',
    padding: '10px',
    borderRadius: '6px',
    fontWeight: '500',
    cursor: 'pointer'
  },
  emptyText: {
    color: '#6b7280',
    fontStyle: 'italic',
    fontSize: '0.9rem'
  },
  eventList: {
    display: 'flex',
    flexDirection: 'column',
    gap: '16px'
  },
  eventItem: {
    padding: '16px',
    border: '1px solid #f3f4f6',
    borderRadius: '8px',
    backgroundColor: '#fdfdfd'
  },
  eventTitle: {
    margin: '0 0 4px 0',
    fontSize: '1.1rem',
    color: '#111827'
  },
  eventDate: {
    margin: '0 0 8px 0',
    fontSize: '0.85rem',
    color: '#3b82f6',
    fontWeight: '500'
  },
  eventDesc: {
    margin: '0 0 12px 0',
    fontSize: '0.9rem',
    color: '#4b5563'
  },
  rsvpSection: {
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
    borderTop: '1px solid #f3f4f6',
    paddingTop: '12px'
  },
  rsvpCounts: {
    display: 'flex',
    gap: '16px',
    fontSize: '0.85rem',
    color: '#6b7280',
    fontWeight: '500'
  },
  rsvpButtons: {
    display: 'flex',
    gap: '8px'
  },
  rsvpButton: {
    flex: 1,
    padding: '6px 0',
    border: '1px solid #d1d5db',
    borderRadius: '6px',
    backgroundColor: '#ffffff',
    color: '#374151',
    cursor: 'pointer',
    fontSize: '0.85rem',
    fontWeight: '500'
  },
  activeGoing: {
    backgroundColor: '#10b981',
    color: '#ffffff',
    borderColor: '#10b981'
  },
  activeNotGoing: {
    backgroundColor: '#ef4444',
    color: '#ffffff',
    borderColor: '#ef4444'
  }
};
