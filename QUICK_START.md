# Quick Start Guide

## 🚀 Get Running in 5 Minutes

### Step 1: Install GCC (Required for SQLite)
```bash
# Using Chocolatey (recommended)
choco install mingw

# Verify
gcc --version
```

### Step 2: Start Backend
```bash
cd social-network/backend
go mod tidy
$env:CGO_ENABLED=1
go run cmd/app/main.go
```

Backend will start on `http://localhost:8080`

### Step 3: Start Frontend (New Terminal)
```bash
cd social-network/frontend
npm install
npm run dev
```

Frontend will start on `http://localhost:3000`

### Step 4: Test the App
1. Open browser → `http://localhost:3000`
2. Click "Register" → Create account
3. Login with credentials
4. Click "My Profile" → View/edit profile
5. Toggle Public/Private
6. Click "Logout"

## ✅ Validation (Optional)

### Backend Validation
```bash
cd social-network/backend
go run validate.go
```

### Frontend Validation
```bash
cd social-network/frontend
node validate.js
```

### Run Unit Tests
```bash
cd social-network/backend
go test -v ./models/...
```

## 🧪 API Testing

### Using PowerShell Script
```bash
cd social-network/backend
.\test_api.ps1
```

### Using cURL
```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","first_name":"John","last_name":"Doe","date_of_birth":"1990-01-01"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"test@example.com","password":"password123"}'

# Get current user
curl http://localhost:8080/api/auth/me -b cookies.txt
```

## 📁 Project Structure

```
social-network/
├── backend/
│   ├── cmd/app/main.go              # Entry point
│   ├── models/user.go                # User & session models
│   ├── pkg/
│   │   ├── db/sqlite/sqlite.go      # Database setup
│   │   └── handlers/                # API handlers
│   ├── queries/middleware/          # Auth & CORS middleware
│   └── server/server.go             # Route configuration
│
└── frontend/
    └── src/
        ├── app/                      # Pages (Next.js)
        ├── context/AuthContext.jsx  # Auth state
        ├── lib/auth.js              # API calls
        └── middleware.js            # Route protection
```

## 🔑 Key Features

✅ User registration with validation  
✅ Secure login (bcrypt + sessions)  
✅ HttpOnly cookies  
✅ Public/private profiles  
✅ Route protection  
✅ CORS configured  
✅ Responsive UI  

## 🆘 Troubleshooting

### "CGO_ENABLED=0" Error
**Solution:** Install GCC/MinGW and set `$env:CGO_ENABLED=1`

### CORS Error in Browser
**Solution:** Backend CORS middleware already configured for localhost:3000

### Port Already in Use
```bash
# Backend (8080)
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Frontend (3000)
netstat -ano | findstr :3000
taskkill /PID <PID> /F
```

### Session Not Persisting
**Solution:** Check browser allows cookies, credentials: 'include' in fetch

## 📚 Documentation

- `SETUP.md` - Detailed setup instructions
- `TESTING.md` - Complete testing guide
- `TEST_RESULTS.md` - Test results and coverage
- `IMPLEMENTATION_SUMMARY.md` - Feature implementation details

## 🎯 What's Implemented

### Backend
- ✅ User registration & login
- ✅ Session management
- ✅ Profile viewing & privacy
- ✅ Authentication middleware
- ✅ CORS configuration
- ✅ SQLite with migrations

### Frontend
- ✅ Registration form
- ✅ Login form
- ✅ Profile page
- ✅ Route protection
- ✅ Auth context
- ✅ Responsive design

## 🔜 Next Steps

Ready to add:
- Followers system
- Posts with privacy levels
- Groups & events
- Real-time chat
- Notifications
