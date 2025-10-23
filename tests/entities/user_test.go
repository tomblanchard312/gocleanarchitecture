package entities_test

import (
	"gocleanarchitecture/entities"
	"strings"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	user, err := entities.NewUser("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error creating user, got %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username to be 'testuser', got %s", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email to be 'test@example.com', got %s", user.Email)
	}

	if user.FullName != "Test User" {
		t.Errorf("Expected full name to be 'Test User', got %s", user.FullName)
	}

	if user.PasswordHash == "" {
		t.Error("Expected password hash to be set")
	}

	if user.PasswordHash == "password123" {
		t.Error("Password should be hashed, not stored in plain text")
	}

	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestNewUserValidation(t *testing.T) {
	testCases := []struct {
		username    string
		email       string
		password    string
		fullName    string
		expectedErr string
	}{
		{"", "test@example.com", "password123", "Test User", "username cannot be empty"},
		{"ab", "test@example.com", "password123", "Test User", "username must be at least 3 characters"},
		{"testuser", "", "password123", "Test User", "email cannot be empty"},
		{"testuser", "invalid-email", "password123", "Test User", "invalid email format"},
		{"testuser", "test@example.com", "pass12", "Test User", "password must be at least 8 characters"},
	}

	for _, tc := range testCases {
		_, err := entities.NewUser(tc.username, tc.email, tc.password, tc.fullName)
		if err == nil {
			t.Errorf("Expected error for username=%s, email=%s, password=%s, fullName=%s, got nil",
				tc.username, tc.email, tc.password, tc.fullName)
			continue
		}

		if !strings.Contains(err.Error(), tc.expectedErr) {
			t.Errorf("Expected error to contain '%s', got '%v'", tc.expectedErr, err)
		}
	}
}

func TestVerifyPassword(t *testing.T) {
	user, err := entities.NewUser("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error creating user, got %v", err)
	}

	// Test correct password
	if !user.VerifyPassword("password123") {
		t.Error("Expected password verification to succeed with correct password")
	}

	// Test incorrect password
	if user.VerifyPassword("wrongpassword") {
		t.Error("Expected password verification to fail with wrong password")
	}
}

func TestUpdate(t *testing.T) {
	user, err := entities.NewUser("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error creating user, got %v", err)
	}

	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Update profile
	err = user.Update("Updated Name", "This is my bio", "http://example.com/avatar.jpg")
	if err != nil {
		t.Fatalf("Expected no error updating user, got %v", err)
	}

	if user.FullName != "Updated Name" {
		t.Errorf("Expected full name to be 'Updated Name', got %s", user.FullName)
	}

	if user.Bio != "This is my bio" {
		t.Errorf("Expected bio to be 'This is my bio', got %s", user.Bio)
	}

	if user.AvatarURL != "http://example.com/avatar.jpg" {
		t.Errorf("Expected avatar URL to be 'http://example.com/avatar.jpg', got %s", user.AvatarURL)
	}

	if !user.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("Expected UpdatedAt (%v) to be after original (%v)", user.UpdatedAt, originalUpdatedAt)
	}
}

func TestUpdateValidation(t *testing.T) {
	user, err := entities.NewUser("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error creating user, got %v", err)
	}

	// Full name too long (>100 characters)
	longName := strings.Repeat("a", 101)
	err = user.Update(longName, "bio", "")
	if err == nil {
		t.Error("Expected error for full name too long, got nil")
	}

	// Bio too long (>500 characters)
	longBio := strings.Repeat("a", 501)
	err = user.Update("Valid Name", longBio, "")
	if err == nil {
		t.Error("Expected error for bio too long, got nil")
	}
}

func TestChangePassword(t *testing.T) {
	user, err := entities.NewUser("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error creating user, got %v", err)
	}

	originalHash := user.PasswordHash

	// Change password
	err = user.ChangePassword("password123", "newpassword456")
	if err != nil {
		t.Fatalf("Expected no error changing password, got %v", err)
	}

	if user.PasswordHash == originalHash {
		t.Error("Expected password hash to be different after change")
	}

	// Verify old password doesn't work
	if user.VerifyPassword("password123") {
		t.Error("Expected old password to not work after change")
	}

	// Verify new password works
	if !user.VerifyPassword("newpassword456") {
		t.Error("Expected new password to work after change")
	}
}

func TestChangePasswordWrongOldPassword(t *testing.T) {
	user, err := entities.NewUser("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error creating user, got %v", err)
	}

	// Try to change password with wrong old password
	err = user.ChangePassword("wrongpassword", "newpassword456")
	if err == nil {
		t.Fatal("Expected error for wrong old password, got nil")
	}

	if err.Error() != "old password is incorrect" {
		t.Errorf("Expected 'old password is incorrect' error, got %v", err)
	}
}

func TestChangePasswordValidation(t *testing.T) {
	user, err := entities.NewUser("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error creating user, got %v", err)
	}

	// New password too short
	err = user.ChangePassword("password123", "pass")
	if err == nil {
		t.Error("Expected error for short new password, got nil")
	}
}

func TestSanitize(t *testing.T) {
	user, err := entities.NewUser("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error creating user, got %v", err)
	}

	if user.PasswordHash == "" {
		t.Error("Expected original user to have password hash")
	}

	sanitized := user.Sanitize()

	if sanitized.PasswordHash != "" {
		t.Error("Expected sanitized user to have empty password hash")
	}

	if sanitized.Username != user.Username {
		t.Error("Expected sanitized user to have same username")
	}

	if sanitized.Email != user.Email {
		t.Error("Expected sanitized user to have same email")
	}
}

