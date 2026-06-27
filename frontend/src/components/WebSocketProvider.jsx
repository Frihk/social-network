'use client';

import { useWebSocket } from '../hooks/useWebSocket';

export default function WebSocketProvider({ children }) {
  useWebSocket({
    onMessage: (message) => {
      if (message?.type === 'notification') {
        window.dispatchEvent(new CustomEvent('realtime-notification', { detail: message }));
      }
    },
  });

  return children;
}
