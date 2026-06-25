'use client';

import { useEffect, useRef, useState, useContext } from 'react';
import { useWebSocket } from '../hooks/useWebSocket';
import { AuthContext } from '../context/AuthContext';
import { getChatHistory, getGroupMessages } from '../lib/chat';
import { EmojiToggleButton } from './EmojiPicker';

export default function ChatWindow({ conversationType, conversationId, conversationName, onClose }) {
  const { user } = useContext(AuthContext);
  const [messages, setMessages] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [messageInput, setMessageInput] = useState('');
  const [error, setError] = useState(null);
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const messagesEndRef = useRef(null);
  const inputRef = useRef(null);
  const pendingMessagesRef = useRef(new Set());
  const { sendMessage, isConnected, lastMessage } = useWebSocket();

  // Format timestamp
  const formatTime = (timestamp) => {
    if (!timestamp) return '';
    const date = new Date(timestamp);
    const now = new Date();
    const isToday = date.toDateString() === now.toDateString();
    
    const timeStr = date.toLocaleTimeString('en-US', { 
      hour: '2-digit', 
      minute: '2-digit',
      hour12: false 
    });
    
    if (isToday) {
      return `Today ${timeStr}`;
    }
    
    const dateStr = date.toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric' 
    });
    return `${dateStr} ${timeStr}`;
  };

  // Auto-scroll to bottom
  const scrollToBottom = () => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  };

  // Fetch message history on mount
  useEffect(() => {
    const fetchHistory = async () => {
      setIsLoading(true);
      setError(null);
      try {
        let history;
        if (conversationType === 'private') {
          history = await getChatHistory(conversationId);
        } else if (conversationType === 'group') {
          history = await getGroupMessages(conversationId);
        }
        
        // Sort by created_at ascending
        const sortedHistory = (history || []).sort((a, b) => {
          const dateA = new Date(a.created_at || 0);
          const dateB = new Date(b.created_at || 0);
          return dateA - dateB;
        });
        
        setMessages(sortedHistory);
      } catch (err) {
        setError(err.message);
      } finally {
        setIsLoading(false);
      }
    };

    fetchHistory();
  }, [conversationType, conversationId]);

  // Handle incoming messages via lastMessage
  useEffect(() => {
    if (!lastMessage || !user) return;

    const msg = lastMessage;
    
    // Filter by conversation type and ID
    if (conversationType === 'private') {
      if (msg.type === 'private_message' && 
          (msg.user_id === conversationId || msg.sender_id === conversationId)) {
        setMessages(prev => {
          // Check if this is replacing an optimistic message
          const optimisticIndex = prev.findIndex(m => 
            m.id?.startsWith('temp-') && 
            m.content === msg.content && 
            m.sender_id === msg.sender_id
          );
          
          if (optimisticIndex >= 0) {
            // Replace optimistic message with real one
            const newMessages = [...prev];
            newMessages[optimisticIndex] = msg;
            pendingMessagesRef.current.delete(prev[optimisticIndex].id);
            return newMessages;
          }
          
          const exists = prev.some(m => m.id === msg.id);
          if (exists) return prev;
          return [...prev, msg].sort((a, b) => {
            const dateA = new Date(a.created_at || 0);
            const dateB = new Date(b.created_at || 0);
            return dateA - dateB;
          });
        });
      }
    } else if (conversationType === 'group') {
      if (msg.type === 'group_message' && msg.group_id === conversationId) {
        setMessages(prev => {
          // Check if this is replacing an optimistic message
          const optimisticIndex = prev.findIndex(m => 
            m.id?.startsWith('temp-') && 
            m.content === msg.content && 
            m.sender_id === msg.sender_id
          );
          
          if (optimisticIndex >= 0) {
            // Replace optimistic message with real one
            const newMessages = [...prev];
            newMessages[optimisticIndex] = msg;
            pendingMessagesRef.current.delete(prev[optimisticIndex].id);
            return newMessages;
          }
          
          const exists = prev.some(m => m.id === msg.id);
          if (exists) return prev;
          return [...prev, msg].sort((a, b) => {
            const dateA = new Date(a.created_at || 0);
            const dateB = new Date(b.created_at || 0);
            return dateA - dateB;
          });
        });
      }
    }
  }, [lastMessage, user, conversationType, conversationId]);

  // Auto-scroll on new messages
  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // Handle emoji select
  const handleEmojiSelect = (emoji) => {
    setMessageInput(prev => prev + emoji);
    // Focus back to input
    if (inputRef.current) {
      inputRef.current.focus();
    }
  };

  // Handle send message
  const handleSend = () => {
    if (!messageInput.trim() || !user) return;

    // Create optimistic message
    const optimisticId = `temp-${Date.now()}`;
    const optimisticMessage = {
      id: optimisticId,
      type: conversationType === 'private' ? 'private_message' : 'group_message',
      content: messageInput.trim(),
      sender_id: user.id,
      user_id: user.id,
      ...(conversationType === 'group' && { group_id: conversationId }),
      created_at: new Date().toISOString(),
    };

    // Add optimistic message to state
    setMessages(prev => [...prev, optimisticMessage].sort((a, b) => {
      const dateA = new Date(a.created_at || 0);
      const dateB = new Date(b.created_at || 0);
      return dateA - dateB;
    }));
    pendingMessagesRef.current.add(optimisticId);

    const message = {
      type: conversationType === 'private' ? 'private_message' : 'group_message',
      content: messageInput.trim(),
      user_id: user.id,
      ...(conversationType === 'group' && { group_id: conversationId }),
    };

    sendMessage(message);
    setMessageInput('');
  };

  // Handle Enter key
  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  if (error) {
    return (
      <div className="flex flex-col h-full bg-white">
        <div className="flex items-center justify-between p-4 border-b">
          <h2 className="text-lg font-semibold">{conversationName}</h2>
          <button 
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700"
          >
            ✕
          </button>
        </div>
        <div className="flex-1 flex items-center justify-center">
          <p className="text-red-500">{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full bg-white">
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b">
        <h2 className="text-lg font-semibold">{conversationName}</h2>
        <button 
          onClick={onClose}
          className="text-gray-500 hover:text-gray-700"
        >
          ✕
        </button>
      </div>

      {/* Connection status */}
      <div className="px-4 py-1 text-xs text-gray-500">
        {isConnected ? (
          <span className="text-green-500">✓ Connected</span>
        ) : (
          <span className="text-yellow-500">Connecting...</span>
        )}
      </div>

      {/* Messages container */}
      <div className="flex-1 overflow-y-auto p-4 space-y-2">
        {isLoading ? (
          <div className="space-y-2">
            {[1, 2, 3].map((i) => (
              <div key={i} className="animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                <div className="h-3 bg-gray-200 rounded w-1/4"></div>
              </div>
            ))}
          </div>
        ) : messages.length === 0 ? (
          <div className="flex items-center justify-center h-full text-gray-500">
            No messages yet. Start the conversation!
          </div>
        ) : (
          messages.map((msg) => {
            const isOwnMessage = msg.sender_id === user?.id || msg.user_id === user?.id;
            return (
              <div
                key={msg.id}
                className={`flex flex-col ${isOwnMessage ? 'items-end' : 'items-start'}`}
              >
                <div
                  className={`max-w-xs px-3 py-2 rounded-lg ${
                    isOwnMessage 
                      ? 'bg-blue-500 text-white' 
                      : 'bg-gray-200 text-black'
                  }`}
                >
                  <p className="text-sm">{msg.content}</p>
                </div>
                <span className="text-xs text-gray-500 mt-1">
                  {formatTime(msg.created_at)}
                </span>
              </div>
            );
          })
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Input area */}
      <div className="p-4 border-t relative">
        <div className="flex gap-2">
          <EmojiToggleButton
            onEmojiSelect={handleEmojiSelect}
            isOpen={showEmojiPicker}
            onToggle={() => setShowEmojiPicker(!showEmojiPicker)}
          />
          <input
            ref={inputRef}
            type="text"
            value={messageInput}
            onChange={(e) => setMessageInput(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Type a message..."
            className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={!isConnected}
          />
          <button
            onClick={handleSend}
            disabled={!isConnected || !messageInput.trim()}
            className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Send
          </button>
        </div>
      </div>
    </div>
  );
}