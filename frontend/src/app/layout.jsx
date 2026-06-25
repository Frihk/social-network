import { AuthProvider } from '@/context/AuthContext';
import { useWebSocket } from '@/hooks/useWebSocket';
import '@/styles/globals.css';

export const metadata = {
  title: 'Social Network',
  description: 'A Facebook-like social network',
};

function WebSocketProvider({ children }) {
  // Initialize WebSocket connection at app root level
  // This ensures a single connection for the entire app
  useWebSocket();

  return children;
}

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          <WebSocketProvider>
            {children}
          </WebSocketProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
