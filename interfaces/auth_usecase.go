package interfaces

import "gocleanarchitecture/entities"

type LoginResponse struct {
	User  *entities.User
	Token string
}

type AuthUseCase interface {
	Register(username, email, password, fullName string) (*LoginResponse, error)
	Login(emailOrUsername, password string) (*LoginResponse, error)
	GetProfile(userID string) (*entities.User, error)
	UpdateProfile(userID, fullName, bio, avatarURL string) (*entities.User, error)
	ChangePassword(userID, oldPassword, newPassword string) error
	GetUserByUsername(username string) (*entities.User, error)
}

