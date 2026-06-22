package models

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			date_of_birth DATE NOT NULL,
			avatar TEXT,
			nickname TEXT,
			about_me TEXT,
			is_private BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)

	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := &User{
		ID:          "test-user-1",
		Email:       "test@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: "1990-01-01",
		IsPrivate:   false,
	}

	err := CreateUser(db, user, "password123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify user was created
	retrievedUser, err := GetUserByEmail(db, user.Email)
	if err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}

	if retrievedUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrievedUser.Email)
	}

	if retrievedUser.FirstName != user.FirstName {
		t.Errorf("Expected first name %s, got %s", user.FirstName, retrievedUser.FirstName)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := &User{
		ID:          "test-user-2",
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Smith",
		DateOfBirth: "1985-05-15",
	}

	CreateUser(db, user, "password123")

	retrievedUser, err := GetUserByEmail(db, user.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrievedUser.ID)
	}
}

func TestGetUserByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := &User{
		ID:          "test-user-3",
		Email:       "jane@example.com",
		FirstName:   "Jane",
		LastName:    "Doe",
		DateOfBirth: "1992-03-20",
	}

	CreateUser(db, user, "password123")

	retrievedUser, err := GetUserByID(db, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if retrievedUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrievedUser.Email)
	}
}

func TestVerifyPassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := &User{
		ID:          "test-user-4",
		Email:       "verify@example.com",
		FirstName:   "Verify",
		LastName:    "Test",
		DateOfBirth: "1990-01-01",
	}

	password := "correctpassword"
	CreateUser(db, user, password)

	retrievedUser, _ := GetUserByEmail(db, user.Email)

	// Test correct password
	err := VerifyPassword(retrievedUser.PasswordHash, password)
	if err != nil {
		t.Errorf("Expected password verification to succeed, got error: %v", err)
	}

	// Test incorrect password
	err = VerifyPassword(retrievedUser.PasswordHash, "wrongpassword")
	if err == nil {
		t.Error("Expected password verification to fail for wrong password")
	}
}

func TestSetProfilePrivacy(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := &User{
		ID:          "test-user-5",
		Email:       "privacy@example.com",
		FirstName:   "Privacy",
		LastName:    "Test",
		DateOfBirth: "1990-01-01",
		IsPrivate:   false,
	}

	CreateUser(db, user, "password123")

	// Set to private
	err := SetProfilePrivacy(db, user.ID, true)
	if err != nil {
		t.Fatalf("Failed to set profile privacy: %v", err)
	}

	// Verify privacy was updated
	retrievedUser, _ := GetUserByID(db, user.ID)
	if !retrievedUser.IsPrivate {
		t.Error("Expected profile to be private")
	}

	// Set back to public
	err = SetProfilePrivacy(db, user.ID, false)
	if err != nil {
		t.Fatalf("Failed to set profile privacy: %v", err)
	}

	retrievedUser, _ = GetUserByID(db, user.ID)
	if retrievedUser.IsPrivate {
		t.Error("Expected profile to be public")
	}
}

func TestSessionManagement(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a user first
	user := &User{
		ID:          "session-user",
		Email:       "session@example.com",
		FirstName:   "Session",
		LastName:    "Test",
		DateOfBirth: "1990-01-01",
	}
	CreateUser(db, user, "password123")

	// Create session
	session := &Session{
		ID:        "test-session-1",
		UserID:    user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err := CreateSession(db, session)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Get session
	retrievedSession, err := GetSessionByID(db, session.ID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if retrievedSession.UserID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, retrievedSession.UserID)
	}

	// Delete session
	err = DeleteSession(db, session.ID)
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	// Verify session is deleted
	_, err = GetSessionByID(db, session.ID)
	if err == nil {
		t.Error("Expected error when getting deleted session")
	}
}
