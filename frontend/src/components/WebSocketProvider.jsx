'use client';

import { useWebSocket } from '../hooks/useWebSocket';

export default function WebSocketProvider({ children }) {
  useWebSocket({
    onMessage: (message) => {
      if (message?.type === 'notification') {
        window.dispatchEvent(new CustomEvent('realtime-notification', { detail: message }));
      } else if (message?.type === 'private_message' || message?.type === 'group_message') {
        window.dispatchEvent(new CustomEvent('realtime-message', { detail: message }));
      }
    },
  });

  return children;
}
