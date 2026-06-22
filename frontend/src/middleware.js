import { NextResponse } from 'next/server';

export function middleware(request) {
  const sessionCookie = request.cookies.get('session_id');
  const isAuthenticated = !!sessionCookie;

  const isAuthPage = request.nextUrl.pathname === '/login' || request.nextUrl.pathname === '/register';

  if (!isAuthenticated && !isAuthPage) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  if (isAuthenticated && isAuthPage) {
    return NextResponse.redirect(new URL('/', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
};
