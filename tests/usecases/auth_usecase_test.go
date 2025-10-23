package usecases_test

import (
	"errors"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/usecases"
	"testing"
)

// Mock Logger
type mockLogger struct{}

func (m *mockLogger) Error(msg string, fields ...interface{}) {}

// Mock UserRepository
type mockUserRepository struct {
	users      map[string]*entities.User
	saveError  error
	findError  error
	existError error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*entities.User),
	}
}

func (m *mockUserRepository) Save(user *entities.User) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) FindByID(id string) (*entities.User, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	return m.users[id], nil
}

func (m *mockUserRepository) FindByEmail(email string) (*entities.User, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) FindByUsername(username string) (*entities.User, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) ExistsByEmail(email string) (bool, error) {
	if m.existError != nil {
		return false, m.existError
	}
	for _, user := range m.users {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockUserRepository) ExistsByUsername(username string) (bool, error) {
	if m.existError != nil {
		return false, m.existError
	}
	for _, user := range m.users {
		if user.Username == username {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockUserRepository) Delete(id string) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepository) GetAll() ([]*entities.User, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	users := make([]*entities.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

// Mock TokenGenerator
type mockTokenGenerator struct {
	token         string
	generateError error
	validateError error
}

func newMockTokenGenerator() *mockTokenGenerator {
	return &mockTokenGenerator{
		token: "mock-jwt-token",
	}
}

func (m *mockTokenGenerator) GenerateToken(userID, username, email string) (string, error) {
	if m.generateError != nil {
		return "", m.generateError
	}
	return m.token, nil
}

func (m *mockTokenGenerator) ValidateToken(token string) (userID string, username string, email string, err error) {
	if m.validateError != nil {
		return "", "", "", m.validateError
	}
	return "user-id", "username", "email@example.com", nil
}

func TestRegister(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockTokenGen := newMockTokenGenerator()
	logger := &mockLogger{}
	authUseCase := usecases.NewAuthUseCase(mockRepo, mockTokenGen, logger)

	response, err := authUseCase.Register("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Token != "mock-jwt-token" {
		t.Errorf("Expected token to be 'mock-jwt-token', got %s", response.Token)
	}

	if response.User.Username != "testuser" {
		t.Errorf("Expected username to be 'testuser', got %s", response.User.Username)
	}

	if response.User.Email != "test@example.com" {
		t.Errorf("Expected email to be 'test@example.com', got %s", response.User.Email)
	}

	if response.User.PasswordHash != "" {
		t.Error("Expected password hash to be sanitized (empty)")
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockTokenGen := newMockTokenGenerator()
	logger := &mockLogger{}
	authUseCase := usecases.NewAuthUseCase(mockRepo, mockTokenGen, logger)

	// Register first user
	_, err := authUseCase.Register("user1", "test@example.com", "password123", "User One")
	if err != nil {
		t.Fatalf("Expected no error on first registration, got %v", err)
	}

	// Try to register with same email
	_, err = authUseCase.Register("user2", "test@example.com", "password456", "User Two")
	if err == nil {
		t.Fatal("Expected error for duplicate email, got nil")
	}

	if err.Error() != "email already registered" {
		t.Errorf("Expected 'email already registered' error, got %v", err)
	}
}

func TestRegisterDuplicateUsername(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockTokenGen := newMockTokenGenerator()
	logger := &mockLogger{}
	authUseCase := usecases.NewAuthUseCase(mockRepo, mockTokenGen, logger)

	// Register first user
	_, err := authUseCase.Register("testuser", "user1@example.com", "password123", "User One")
	if err != nil {
		t.Fatalf("Expected no error on first registration, got %v", err)
	}

	// Try to register with same username
	_, err = authUseCase.Register("testuser", "user2@example.com", "password456", "User Two")
	if err == nil {
		t.Fatal("Expected error for duplicate username, got nil")
	}

	if err.Error() != "username already taken" {
		t.Errorf("Expected 'username already taken' error, got %v", err)
	}
}

func TestLogin(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockTokenGen := newMockTokenGenerator()
	logger := &mockLogger{}
	authUseCase := usecases.NewAuthUseCase(mockRepo, mockTokenGen, logger)

	// Register a user first
	_, err := authUseCase.Register("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error during registration, got %v", err)
	}

	// Test login with email
	response, err := authUseCase.Login("test@example.com", "password123")
	if err != nil {
		t.Fatalf("Expected no error during login with email, got %v", err)
	}

	if response.Token != "mock-jwt-token" {
		t.Errorf("Expected token to be 'mock-jwt-token', got %s", response.Token)
	}

	// Test login with username
	response, err = authUseCase.Login("testuser", "password123")
	if err != nil {
		t.Fatalf("Expected no error during login with username, got %v", err)
	}

	if response.User.Username != "testuser" {
		t.Errorf("Expected username to be 'testuser', got %s", response.User.Username)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockTokenGen := newMockTokenGenerator()
	logger := &mockLogger{}
	authUseCase := usecases.NewAuthUseCase(mockRepo, mockTokenGen, logger)

	// Register a user
	_, err := authUseCase.Register("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error during registration, got %v", err)
	}

	// Test with wrong password
	_, err = authUseCase.Login("test@example.com", "wrongpassword")
	if err == nil {
		t.Fatal("Expected error for wrong password, got nil")
	}

	// Test with non-existent user
	_, err = authUseCase.Login("nonexistent@example.com", "password123")
	if err == nil {
		t.Fatal("Expected error for non-existent user, got nil")
	}
}

func TestGetProfile(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockTokenGen := newMockTokenGenerator()
	logger := &mockLogger{}
	authUseCase := usecases.NewAuthUseCase(mockRepo, mockTokenGen, logger)

	// Register a user
	response, err := authUseCase.Register("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error during registration, got %v", err)
	}

	userID := response.User.ID

	// Get profile
	user, err := authUseCase.GetProfile(userID)
	if err != nil {
		t.Fatalf("Expected no error getting profile, got %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username to be 'testuser', got %s", user.Username)
	}

	if user.PasswordHash != "" {
		t.Error("Expected password hash to be sanitized")
	}
}

func TestUpdateProfile(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockTokenGen := newMockTokenGenerator()
	logger := &mockLogger{}
	authUseCase := usecases.NewAuthUseCase(mockRepo, mockTokenGen, logger)

	// Register a user
	response, err := authUseCase.Register("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error during registration, got %v", err)
	}

	userID := response.User.ID

	// Update profile
	updatedUser, err := authUseCase.UpdateProfile(userID, "Updated Name", "New bio", "http://avatar.com/new.jpg")
	if err != nil {
		t.Fatalf("Expected no error updating profile, got %v", err)
	}

	if updatedUser.FullName != "Updated Name" {
		t.Errorf("Expected full name to be 'Updated Name', got %s", updatedUser.FullName)
	}

	if updatedUser.Bio != "New bio" {
		t.Errorf("Expected bio to be 'New bio', got %s", updatedUser.Bio)
	}

	if updatedUser.AvatarURL != "http://avatar.com/new.jpg" {
		t.Errorf("Expected avatar URL to be 'http://avatar.com/new.jpg', got %s", updatedUser.AvatarURL)
	}
}

func TestChangePassword(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockTokenGen := newMockTokenGenerator()
	logger := &mockLogger{}
	authUseCase := usecases.NewAuthUseCase(mockRepo, mockTokenGen, logger)

	// Register a user
	response, err := authUseCase.Register("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error during registration, got %v", err)
	}

	userID := response.User.ID

	// Change password
	err = authUseCase.ChangePassword(userID, "password123", "newpassword456")
	if err != nil {
		t.Fatalf("Expected no error changing password, got %v", err)
	}

	// Verify new password works
	_, err = authUseCase.Login("test@example.com", "newpassword456")
	if err != nil {
		t.Fatalf("Expected login to work with new password, got %v", err)
	}

	// Verify old password doesn't work
	_, err = authUseCase.Login("test@example.com", "password123")
	if err == nil {
		t.Fatal("Expected error logging in with old password, got nil")
	}
}

func TestChangePasswordWrongOldPassword(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockTokenGen := newMockTokenGenerator()
	logger := &mockLogger{}
	authUseCase := usecases.NewAuthUseCase(mockRepo, mockTokenGen, logger)

	// Register a user
	response, err := authUseCase.Register("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error during registration, got %v", err)
	}

	userID := response.User.ID

	// Try to change password with wrong old password
	err = authUseCase.ChangePassword(userID, "wrongpassword", "newpassword456")
	if err == nil {
		t.Fatal("Expected error for wrong old password, got nil")
	}
}

func TestGetUserByUsername(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockTokenGen := newMockTokenGenerator()
	logger := &mockLogger{}
	authUseCase := usecases.NewAuthUseCase(mockRepo, mockTokenGen, logger)

	// Register a user
	_, err := authUseCase.Register("testuser", "test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Expected no error during registration, got %v", err)
	}

	// Get user by username
	user, err := authUseCase.GetUserByUsername("testuser")
	if err != nil {
		t.Fatalf("Expected no error getting user by username, got %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username to be 'testuser', got %s", user.Username)
	}

	if user.PasswordHash != "" {
		t.Error("Expected password hash to be sanitized")
	}
}

func TestRegisterTokenGenerationError(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockTokenGen := newMockTokenGenerator()
	mockTokenGen.generateError = errors.New("token generation failed")
	logger := &mockLogger{}
	authUseCase := usecases.NewAuthUseCase(mockRepo, mockTokenGen, logger)

	_, err := authUseCase.Register("testuser", "test@example.com", "password123", "Test User")
	if err == nil {
		t.Fatal("Expected error when token generation fails, got nil")
	}

	if err.Error() != "failed to generate authentication token" {
		t.Errorf("Expected 'failed to generate authentication token' error, got %v", err)
	}
}
