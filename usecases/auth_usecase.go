package usecases

import (
	"errors"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"

	"github.com/google/uuid"
)

type TokenGenerator interface {
	GenerateToken(userID, username, email string) (string, error)
	ValidateToken(token string) (userID string, username string, email string, err error)
}

type AuthUseCase struct {
	UserRepo       interfaces.UserRepository
	TokenGenerator TokenGenerator
	Logger         Logger
}

func NewAuthUseCase(userRepo interfaces.UserRepository, tokenGen TokenGenerator, logger Logger) interfaces.AuthUseCase {
	return &AuthUseCase{
		UserRepo:       userRepo,
		TokenGenerator: tokenGen,
		Logger:         logger,
	}
}

// Register creates a new user account
func (u *AuthUseCase) Register(username, email, password, fullName string) (*interfaces.LoginResponse, error) {
	// Check if email already exists
	existingEmail, err := u.UserRepo.ExistsByEmail(email)
	if err != nil {
		u.Logger.Error("Failed to check email existence", "error", err)
		return nil, errors.New("failed to check email availability")
	}
	if existingEmail {
		return nil, errors.New("email already registered")
	}

	// Check if username already exists
	existingUsername, err := u.UserRepo.ExistsByUsername(username)
	if err != nil {
		u.Logger.Error("Failed to check username existence", "error", err)
		return nil, errors.New("failed to check username availability")
	}
	if existingUsername {
		return nil, errors.New("username already taken")
	}

	// Create new user using domain factory
	user, err := entities.NewUser(username, email, password, fullName)
	if err != nil {
		return nil, err
	}

	// Generate unique ID
	user.ID = uuid.New().String()

	// Save user to repository
	err = u.UserRepo.Save(user)
	if err != nil {
		u.Logger.Error("Failed to save user", "error", err)
		return nil, errors.New("failed to create user account")
	}

	// Generate JWT token
	token, err := u.TokenGenerator.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		u.Logger.Error("Failed to generate token", "error", err)
		return nil, errors.New("failed to generate authentication token")
	}

	return &interfaces.LoginResponse{
		User:  user.Sanitize(),
		Token: token,
	}, nil
}

// Login authenticates a user and returns a token
func (u *AuthUseCase) Login(emailOrUsername, password string) (*interfaces.LoginResponse, error) {
	// Try to find user by email first, then by username
	var user *entities.User
	var err error

	user, err = u.UserRepo.FindByEmail(emailOrUsername)
	if err != nil {
		u.Logger.Error("Failed to find user by email", "error", err)
		return nil, errors.New("authentication failed")
	}

	if user == nil {
		user, err = u.UserRepo.FindByUsername(emailOrUsername)
		if err != nil {
			u.Logger.Error("Failed to find user by username", "error", err)
			return nil, errors.New("authentication failed")
		}
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if !user.VerifyPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := u.TokenGenerator.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		u.Logger.Error("Failed to generate token", "error", err)
		return nil, errors.New("failed to generate authentication token")
	}

	return &interfaces.LoginResponse{
		User:  user.Sanitize(),
		Token: token,
	}, nil
}

// GetProfile retrieves a user's profile by ID
func (u *AuthUseCase) GetProfile(userID string) (*entities.User, error) {
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		u.Logger.Error("Failed to get user profile", "error", err, "userID", userID)
		return nil, errors.New("failed to retrieve profile")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user.Sanitize(), nil
}

// UpdateProfile updates a user's profile information
func (u *AuthUseCase) UpdateProfile(userID, fullName, bio, avatarURL string) (*entities.User, error) {
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		u.Logger.Error("Failed to find user for update", "error", err, "userID", userID)
		return nil, errors.New("failed to retrieve user")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Use domain method to update
	err = user.Update(fullName, bio, avatarURL)
	if err != nil {
		return nil, err
	}

	// Save updated user
	err = u.UserRepo.Save(user)
	if err != nil {
		u.Logger.Error("Failed to update user profile", "error", err, "userID", userID)
		return nil, errors.New("failed to update profile")
	}

	return user.Sanitize(), nil
}

// ChangePassword allows a user to change their password
func (u *AuthUseCase) ChangePassword(userID, oldPassword, newPassword string) error {
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		u.Logger.Error("Failed to find user for password change", "error", err, "userID", userID)
		return errors.New("failed to retrieve user")
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Use domain method to change password
	err = user.ChangePassword(oldPassword, newPassword)
	if err != nil {
		return err
	}

	// Save updated user
	err = u.UserRepo.Save(user)
	if err != nil {
		u.Logger.Error("Failed to save password change", "error", err, "userID", userID)
		return errors.New("failed to change password")
	}

	return nil
}

// GetUserByUsername retrieves a user by username (for public profiles)
func (u *AuthUseCase) GetUserByUsername(username string) (*entities.User, error) {
	user, err := u.UserRepo.FindByUsername(username)
	if err != nil {
		u.Logger.Error("Failed to find user by username", "error", err, "username", username)
		return nil, errors.New("failed to retrieve user")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user.Sanitize(), nil
}
