package entities

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	FullName     string
	Bio          string
	AvatarURL    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
)

// NewUser creates a new user with validation and password hashing
func NewUser(username, email, password, fullName string) (*User, error) {
	if err := validateUserData(username, email, password); err != nil {
		return nil, err
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	now := time.Now()
	return &User{
		Username:     strings.TrimSpace(username),
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: passwordHash,
		FullName:     strings.TrimSpace(fullName),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Update updates user profile information
func (u *User) Update(fullName, bio, avatarURL string) error {
	if len(fullName) > 100 {
		return errors.New("full name cannot exceed 100 characters")
	}
	if len(bio) > 500 {
		return errors.New("bio cannot exceed 500 characters")
	}

	u.FullName = strings.TrimSpace(fullName)
	u.Bio = strings.TrimSpace(bio)
	u.AvatarURL = strings.TrimSpace(avatarURL)
	u.UpdatedAt = time.Now()
	return nil
}

// VerifyPassword checks if the provided password matches the user's password
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// ChangePassword updates the user's password
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	if !u.VerifyPassword(oldPassword) {
		return errors.New("old password is incorrect")
	}

	if err := validatePassword(newPassword); err != nil {
		return err
	}

	passwordHash, err := hashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	u.PasswordHash = passwordHash
	u.UpdatedAt = time.Now()
	return nil
}

// Sanitize removes sensitive information before returning to client
func (u *User) Sanitize() *User {
	return &User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Bio:       u.Bio,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		// PasswordHash is intentionally omitted
	}
}

// Private helper functions

func validateUserData(username, email, password string) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	if err := validateEmail(email); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	return nil
}

func validateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(username) > 30 {
		return errors.New("username cannot exceed 30 characters")
	}
	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}
	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if len(password) > 72 {
		return errors.New("password cannot exceed 72 characters")
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

