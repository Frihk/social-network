'use client';

import { useRef, useState, useEffect } from 'react';

// Common emojis organized by category
const EMOJI_CATEGORIES = {
  smileys: ['😀', '😁', '😂', '🤣', '😃', '😄', '😅', '😆', '😉', '😊', '🙂', '🤗', '🤩', '😍', '😘', '😚', '😙', '😋', '😛', '😜', '🤪', '😌', '😔', '😒', '😲', '😟', '😕', '🙁', '☹️', '😮', '😯', '😨', '😰', '😥', '😢', '😭', '😱', '😖', '😣', '😞', '😓', '😩', '😫', '🥱', '😤', '😡', '😠', '🤬', '😈', '👿', '💀', '☠️'],
  hearts: ['❤️', '🧡', '💛', '💚', '💙', '💜', '🖤', '🤍', '🤎', '💔', '💕', '💞', '💓', '💗', '💖', '💘', '💝', '💟'],
  hands: ['👋', '🤚', '🖐️', '✋', '🖖', '👌', '🤌', '🤏', '✌️', '🤞', '🫰', '🤟', '🤘', '🤙', '👍', '👎', '👊', '👏', '🙌', '👐'],
  other: ['👀', '👁️', '👅', '⭐', '✨', '💫', '⚡', '🔥', '🎉', '🎊', '🎈', '🎀', '🎁', '🍕', '🍔', '🍟', '🍗', '🍖', '🌮', '🌯', '🍿', '☕', '🍵', '🍷', '🍾', '🍻'],
};

// Flatten all emojis into a single array
const ALL_EMOJIS = [...EMOJI_CATEGORIES.smileys, ...EMOJI_CATEGORIES.hearts, ...EMOJI_CATEGORIES.hands, ...EMOJI_CATEGORIES.other];

export default function EmojiPicker({ onEmojiSelect, isOpen, onToggle }) {
  const pickerRef = useRef(null);
  const [searchQuery, setSearchQuery] = useState('');

  // Close picker when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (pickerRef.current && !pickerRef.current.contains(event.target)) {
        onToggle();
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen, onToggle]);

  // Handle Escape key
  useEffect(() => {
    const handleEscape = (event) => {
      if (event.key === 'Escape' && isOpen) {
        onToggle();
      }
    };

    if (isOpen) {
      document.addEventListener('keydown', handleEscape);
    }

    return () => {
      document.removeEventListener('keydown', handleEscape);
    };
  }, [isOpen, onToggle]);

  // Filter emojis by search query (simple substring match)
  const filteredEmojis = searchQuery
    ? ALL_EMOJIS.filter(emoji => {
        // This is a simple filter - in a real app you'd have emoji names/keywords
        return true; // For now, just return all emojis
      })
    : ALL_EMOJIS;

  if (!isOpen) return null;

  return (
    <div
      ref={pickerRef}
      className="absolute bottom-full mb-2 bg-white border border-gray-200 rounded-lg shadow-lg p-2 w-80 max-h-96 overflow-y-auto z-10"
      role="dialog"
      aria-label="Emoji picker"
    >
      {/* Search bar (optional, for future enhancement) */}
      <div className="mb-2">
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search emojis..."
          className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
        />
      </div>

      {/* Emoji grid */}
      <div className="grid grid-cols-8 gap-1">
        {filteredEmojis.map((emoji, index) => (
          <button
            key={index}
            onClick={() => {
              onEmojiSelect(emoji);
              onToggle();
              setSearchQuery('');
            }}
            className="w-8 h-8 flex items-center justify-center text-lg hover:bg-gray-100 rounded transition-transform hover:scale-110 focus:outline-none focus:ring-1 focus:ring-blue-500"
            aria-label={`Select emoji ${emoji}`}
            type="button"
          >
            {emoji}
          </button>
        ))}
      </div>
    </div>
  );
}

// Toggle button component
export function EmojiToggleButton({ onEmojiSelect, isOpen, onToggle }) {
  return (
    <div className="relative">
      <button
        onClick={onToggle}
        className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded"
        aria-label="Toggle emoji picker"
        aria-expanded={isOpen}
        type="button"
      >
        😀
      </button>
      <EmojiPicker
        onEmojiSelect={onEmojiSelect}
        isOpen={isOpen}
        onToggle={onToggle}
      />
    </div>
  );
}