'use client';

import { useContext, useEffect, useRef, useState, useCallback } from 'react';
import { AuthContext } from '../context/AuthContext';

const WS_URL = 'ws://localhost:8080/api/ws';

export function useWebSocket({ onMessage, maxRetries = 5, retryDelay = 3000 } = {}) {
  const { sessionToken } = useContext(AuthContext);
  const [isConnected, setIsConnected] = useState(false);
  const [lastMessage, setLastMessage] = useState(null);
  const wsRef = useRef(null);
  const messageHandlersRef = useRef([]);
  const reconnectAttemptsRef = useRef(0);
  const reconnectTimeoutRef = useRef(null);
  const messageQueueRef = useRef([]);
  const handlerAddedRef = useRef(false);

  // Register message handler
  useEffect(() => {
    if (onMessage && typeof onMessage === 'function' && !handlerAddedRef.current) {
      messageHandlersRef.current.push(onMessage);
      handlerAddedRef.current = true;
    }
    return () => {
      if (onMessage && typeof onMessage === 'function') {
        messageHandlersRef.current = messageHandlersRef.current.filter(
          (handler) => handler !== onMessage
        );
        handlerAddedRef.current = false;
      }
    };
  }, [onMessage]);

  // Connect to WebSocket when sessionToken is available
  useEffect(() => {
    if (!sessionToken) {
      return;
    }

    let isUnmounted = false;

    const connect = () => {
      if (isUnmounted) return;

      const wsUrl = `${WS_URL}?session_token=${sessionToken}`;
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        if (isUnmounted) return;
        console.log('WebSocket connected');
        setIsConnected(true);
        reconnectAttemptsRef.current = 0;

        // Send ping to confirm connection
        sendMessage({ type: 'ping' });

        // Send any queued messages
        while (messageQueueRef.current.length > 0) {
          const queuedMessage = messageQueueRef.current.shift();
          if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify(queuedMessage));
          }
        }
      };

      ws.onmessage = (event) => {
        if (isUnmounted) return;

        try {
          const message = JSON.parse(event.data);
          setLastMessage(message);

          // Call all registered handlers
          messageHandlersRef.current.forEach((handler) => {
            try {
              handler(message);
            } catch (handlerError) {
              console.error('Error in message handler:', handlerError);
            }
          });
        } catch (parseError) {
          console.error('Error parsing WebSocket message:', parseError);
        }
      };

      ws.onclose = (event) => {
        if (isUnmounted) return;
        console.log('WebSocket disconnected:', event.code, event.reason);
        setIsConnected(false);

        // Attempt to reconnect if not at max retries
        if (reconnectAttemptsRef.current < maxRetries) {
          reconnectAttemptsRef.current += 1;
          const delay = retryDelay * Math.pow(2, reconnectAttemptsRef.current - 1);
          console.log(`Attempting to reconnect (${reconnectAttemptsRef.current}/${maxRetries}) in ${delay}ms`);

          reconnectTimeoutRef.current = setTimeout(connect, delay);
        } else {
          console.log('Max reconnection attempts reached');
        }
      };

      ws.onerror = (error) => {
        if (isUnmounted) return;
        console.error('WebSocket error:', error);
      };
    };

    connect();

    return () => {
      isUnmounted = true;
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [sessionToken, maxRetries, retryDelay]);

  // Send message function
  const sendMessage = useCallback((payload) => {
    if (!payload || typeof payload !== 'object') {
      console.error('sendMessage: payload must be an object');
      return;
    }

    const message = {
      type: payload.type || 'message',
      user_id: payload.user_id || null,
      group_id: payload.group_id || null,
      content: payload.content || '',
      data: payload.data || null,
    };

    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    } else {
      // Queue message for later
      messageQueueRef.current.push(message);
    }
  }, []);

  return {
    sendMessage,
    isConnected,
    lastMessage,
  };
}