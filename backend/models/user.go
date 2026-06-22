package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	DateOfBirth  string    `json:"date_of_birth"`
	Avatar       *string   `json:"avatar"`
	Nickname     *string   `json:"nickname"`
	AboutMe      *string   `json:"about_me"`
	IsPrivate    bool      `json:"is_private"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func CreateUser(db *sql.DB, user *User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)

	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_private)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = db.Exec(query, user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.DateOfBirth, user.Avatar, user.Nickname, user.AboutMe, user.IsPrivate)

	return err
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, first_name, last_name, date_of_birth, 
		       avatar, nickname, about_me, is_private, created_at, updated_at
		FROM users WHERE email = ?
	`

	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsPrivate,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return user, err
}

func GetUserByID(db *sql.DB, id string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, email, password_hash, first_name, last_name, date_of_birth, 
		       avatar, nickname, about_me, is_private, created_at, updated_at
		FROM users WHERE id = ?
	`

	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsPrivate,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return user, err
}

func UpdateUser(db *sql.DB, user *User) error {
	query := `
		UPDATE users 
		SET first_name = ?, last_name = ?, avatar = ?, nickname = ?, about_me = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := db.Exec(query, user.FirstName, user.LastName, user.Avatar, user.Nickname, user.AboutMe, user.ID)
	return err
}

func SetProfilePrivacy(db *sql.DB, userID string, isPrivate bool) error {
	query := `UPDATE users SET is_private = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, isPrivate, userID)
	return err
}

func CreateSession(db *sql.DB, session *Session) error {
	query := `INSERT INTO sessions (id, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, session.ID, session.UserID, session.CreatedAt, session.ExpiresAt)
	return err
}

func GetSessionByID(db *sql.DB, sessionID string) (*Session, error) {
	session := &Session{}
	query := `SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = ?`

	err := db.QueryRow(query, sessionID).Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("session not found")
	}

	return session, err
}

func DeleteSession(db *sql.DB, sessionID string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := db.Exec(query, sessionID)
	return err
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
