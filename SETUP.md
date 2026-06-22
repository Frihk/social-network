# Social Network Setup Guide

## Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd social-network/backend
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run cmd/app/main.go
   ```

   The backend will start on `http://localhost:8080`

## Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd social-network/frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Run the development server:
   ```bash
   npm run dev
   ```

   The frontend will start on `http://localhost:3000`

## Testing the Application

1. Open your browser and go to `http://localhost:3000`
2. You'll be redirected to the login page
3. Click "Register" to create a new account
4. Fill in the registration form with:
   - Email
   - Password
   - First Name
   - Last Name
   - Date of Birth
   - (Optional) Nickname, About Me, Avatar URL
5. After registration, you'll be redirected to the login page
6. Login with your credentials
7. You'll be redirected to the home page
8. Click "My Profile" to view your profile
9. Toggle between Public/Private profile
10. Click "Logout" to end your session

## API Endpoints

### Public Endpoints
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login

### Protected Endpoints (require session cookie)
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user
- `GET /api/users/:id` - Get user profile
- `PUT /api/users/:id/privacy` - Toggle profile privacy

## Database

The application uses SQLite with the database file stored at:
`social-network/backend/social-network.db`

Migrations are automatically applied on startup from:
`social-network/backend/pkg/db/migration/sqlite/`
