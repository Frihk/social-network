<<<<<<< HEAD
//
=======
import { AuthProvider } from '@/context/AuthContext';
import '@/styles/globals.css';

export const metadata = {
  title: 'Social Network',
  description: 'A Facebook-like social network',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}
>>>>>>> kagunda
