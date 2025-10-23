package db

import (
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
	"strings"
	"sync"
)

type InMemoryUserRepository struct {
	users map[string]*entities.User
	mu    sync.RWMutex
}

func NewInMemoryUserRepository() interfaces.UserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*entities.User),
	}
}

func (r *InMemoryUserRepository) Save(user *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create a copy to avoid external modifications
	userCopy := *user
	r.users[user.ID] = &userCopy
	return nil
}

func (r *InMemoryUserRepository) FindByID(id string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}

	// Return a copy to avoid external modifications
	userCopy := *user
	return &userCopy, nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	email = strings.ToLower(strings.TrimSpace(email))
	for _, user := range r.users {
		if strings.ToLower(user.Email) == email {
			// Return a copy to avoid external modifications
			userCopy := *user
			return &userCopy, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) FindByUsername(username string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	username = strings.ToLower(strings.TrimSpace(username))
	for _, user := range r.users {
		if strings.ToLower(user.Username) == username {
			// Return a copy to avoid external modifications
			userCopy := *user
			return &userCopy, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) ExistsByEmail(email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	email = strings.ToLower(strings.TrimSpace(email))
	for _, user := range r.users {
		if strings.ToLower(user.Email) == email {
			return true, nil
		}
	}
	return false, nil
}

func (r *InMemoryUserRepository) ExistsByUsername(username string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	username = strings.ToLower(strings.TrimSpace(username))
	for _, user := range r.users {
		if strings.ToLower(user.Username) == username {
			return true, nil
		}
	}
	return false, nil
}

func (r *InMemoryUserRepository) GetAll() ([]*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*entities.User, 0, len(r.users))
	for _, user := range r.users {
		// Return a copy to avoid external modifications
		userCopy := *user
		users = append(users, &userCopy)
	}

	return users, nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.users, id)
	return nil
}
