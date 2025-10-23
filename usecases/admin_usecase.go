package usecases

import (
	"errors"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
)

type AdminUseCase struct {
	UserRepo interfaces.UserRepository
	Logger   Logger
}

func NewAdminUseCase(userRepo interfaces.UserRepository, logger Logger) *AdminUseCase {
	return &AdminUseCase{
		UserRepo: userRepo,
		Logger:   logger,
	}
}

// GetAllUsers retrieves all users from the system
func (uc *AdminUseCase) GetAllUsers() ([]*entities.User, error) {
	users, err := uc.UserRepo.GetAll()
	if err != nil {
		uc.Logger.Error("Admin: Failed to fetch all users", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("failed to fetch users")
	}

	return users, nil
}

// GetUserByID retrieves a user by their ID
func (uc *AdminUseCase) GetUserByID(userID string) (*entities.User, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := uc.UserRepo.FindByID(userID)
	if err != nil {
		uc.Logger.Error("Admin: Failed to fetch user", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		return nil, errors.New("failed to fetch user")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UpdateUserRole updates a user's role
func (uc *AdminUseCase) UpdateUserRole(userID string, newRole entities.UserRole) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Fetch user
	user, err := uc.UserRepo.FindByID(userID)
	if err != nil {
		uc.Logger.Error("Admin: Failed to fetch user", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		return errors.New("failed to fetch user")
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Update role
	if err := user.SetRole(newRole); err != nil {
		return err
	}

	// Save updated user
	if err := uc.UserRepo.Save(user); err != nil {
		uc.Logger.Error("Admin: Failed to save user", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		return errors.New("failed to update user role")
	}

	return nil
}

// DeleteUser deletes a user from the system
func (uc *AdminUseCase) DeleteUser(userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Check if user exists
	user, err := uc.UserRepo.FindByID(userID)
	if err != nil {
		uc.Logger.Error("Admin: Failed to fetch user", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		return errors.New("failed to fetch user")
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Delete the user
	if err := uc.UserRepo.Delete(userID); err != nil {
		uc.Logger.Error("Admin: Failed to delete user", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		return errors.New("failed to delete user")
	}

	return nil
}
