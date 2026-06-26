'use client';

import { useWebSocket } from '../hooks/useWebSocket';

export default function WebSocketProvider({ children }) {
  // Initialize WebSocket connection at app root level
  // This ensures a single connection for the entire app
  useWebSocket();

  return children;
}
