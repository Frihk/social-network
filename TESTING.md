# Testing Guide

## Prerequisites

### Windows - Install GCC for SQLite Support

SQLite requires CGO which needs a C compiler. Install MinGW-w64:

1. **Option 1: Using Chocolatey** (recommended)
   ```powershell
   choco install mingw
   ```

2. **Option 2: Manual Installation**
   - Download from: https://sourceforge.net/projects/mingw-w64/
   - Add `C:\mingw64\bin` to your PATH

3. **Verify Installation**
   ```bash
   gcc --version
   ```

## Backend Tests

### Unit Tests

Run the Go unit tests:

```bash
cd social-network/backend
go test -v ./models/...
```

### Test Coverage

```bash
go test -cover ./models/...
```

### What the Tests Cover

#### User Model Tests (`models/user_test.go`)

1. **TestCreateUser** - Tests user creation with password hashing
2. **TestGetUserByEmail** - Tests retrieving users by email
3. **TestGetUserByID** - Tests retrieving users by ID
4. **TestVerifyPassword** - Tests password verification (correct/incorrect)
5. **TestSetProfilePrivacy** - Tests toggling profile privacy
6. **TestSessionManagement** - Tests session creation, retrieval, and deletion

### Integration Tests (Requires Running Server)

Once GCC is installed and the server is running:

```bash
cd social-network/backend
go run cmd/app/main.go
```

In another terminal, run the API tests:

```powershell
cd social-network/backend
.\test_api.ps1
```

### Manual API Testing with cURL

#### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "date_of_birth": "1990-01-01",
    "nickname": "johndoe",
    "about_me": "Hello, I am John!"
  }'
```

#### 2. Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### 3. Get Current User
```bash
curl -X GET http://localhost:8080/api/auth/me \
  -b cookies.txt
```

#### 4. Get User Profile
```bash
curl -X GET http://localhost:8080/api/users/{user_id} \
  -b cookies.txt
```

#### 5. Update Profile Privacy
```bash
curl -X PUT http://localhost:8080/api/users/{user_id}/privacy \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"is_private": true}'
```

#### 6. Logout
```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -b cookies.txt
```

## Frontend Testing

### Install Dependencies

```bash
cd social-network/frontend
npm install
```

### Start Development Server

```bash
npm run dev
```

Visit `http://localhost:3000`

### Manual Testing Flow

1. **Registration Flow**
   - Navigate to `/register`
   - Fill in all required fields
   - Submit form
   - Should redirect to `/login`

2. **Login Flow**
   - Enter registered credentials
   - Click "Login"
   - Should redirect to home page `/`

3. **Authentication Check**
   - Try accessing `/` without logging in
   - Should redirect to `/login`
   - After login, try accessing `/login`
   - Should redirect to `/`

4. **Profile Features**
   - Click "My Profile"
   - View profile information
   - Toggle "Public/Private" button
   - Verify button text changes

5. **Logout Flow**
   - Click "Logout"
   - Should redirect to `/login`
   - Session should be cleared

### Browser DevTools Testing

Open browser DevTools (F12) and check:

1. **Network Tab**
   - API calls to `http://localhost:8080/api/*`
   - Status codes (201, 200, 401, etc.)
   - Request/response payloads

2. **Application Tab > Cookies**
   - `session_id` cookie should be set after login
   - `HttpOnly` flag should be enabled
   - Cookie should be cleared after logout

3. **Console Tab**
   - No JavaScript errors
   - Check for any warnings

## Test Results Expected

### Successful Registration
- Status: 201 Created
- Response: `{"message": "User created successfully"}`

### Successful Login
- Status: 200 OK
- Response: User object with all fields
- Cookie: `session_id` set with HttpOnly flag

### Get Current User (Authenticated)
- Status: 200 OK
- Response: User object

### Get Current User (Unauthenticated)
- Status: 401 Unauthorized
- Response: `{"error": "Unauthorized"}`

### Profile Privacy Toggle
- Status: 200 OK
- Response: `{"is_private": true}` or `{"is_private": false}`

### Logout
- Status: 200 OK
- Response: `{"message": "Logged out successfully"}`
- Cookie: `session_id` cleared

## Common Issues

### Issue: `Binary was compiled with 'CGO_ENABLED=0'`
**Solution:** Install GCC/MinGW-w64 and rebuild with CGO enabled:
```bash
$env:CGO_ENABLED=1
go build ./cmd/app
```

### Issue: CORS errors in browser
**Solution:** Ensure backend CORS middleware allows `http://localhost:3000`

### Issue: 401 Unauthorized on protected routes
**Solution:** 
- Check that session cookie is being sent
- Verify cookie has not expired
- Ensure `credentials: 'include'` is set in fetch calls

### Issue: Frontend not connecting to backend
**Solution:**
- Verify backend is running on port 8080
- Check `API_URL` in `src/lib/auth.js`
- Ensure no firewall blocking localhost connections

## Performance Testing

### Load Testing (Optional)

If you have Apache Bench installed:

```bash
# Test registration endpoint
ab -n 100 -c 10 -p register.json -T application/json http://localhost:8080/api/auth/register

# Test login endpoint  
ab -n 100 -c 10 -p login.json -T application/json http://localhost:8080/api/auth/login
```

## Security Testing Checklist

- [x] Passwords are hashed with bcrypt
- [x] Sessions use secure random UUIDs
- [x] Cookies are HttpOnly
- [x] CORS is configured properly
- [x] SQL injection protection (using parameterized queries)
- [x] Authentication required for protected routes
- [x] Authorization checks (users can only modify their own data)

## Next Steps

After basic authentication tests pass:
1. Test follower system
2. Test posts creation
3. Test groups functionality
4. Test real-time chat
5. Test notifications
