import { NextResponse } from 'next/server';

const protectedPaths = ['/chat', '/create-post', '/groups', '/profile'];

function isProtected(pathname) {
  return protectedPaths.some(p => pathname === p || pathname.startsWith(p + '/'));
}

export function middleware(request) {
  const sessionCookie = request.cookies.get('session_id');
  const isAuthenticated = !!sessionCookie;

  const { pathname } = request.nextUrl;
  const isAuthPage = pathname === '/login' || pathname === '/register';

  if (isAuthenticated && isAuthPage) {
    return NextResponse.redirect(new URL('/', request.url));
  }

  if (!isAuthenticated && isProtected(pathname)) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
};
