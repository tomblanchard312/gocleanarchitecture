package interfaces

import "gocleanarchitecture/entities"

type UserRepository interface {
	Save(user *entities.User) error
	FindByID(id string) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	FindByUsername(username string) (*entities.User, error)
	Delete(id string) error
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
}

