package sqlite

import (
	"database/sql"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteUserRepository struct {
	DB *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) interfaces.UserRepository {
	return &SQLiteUserRepository{DB: db}
}

func (r *SQLiteUserRepository) Save(user *entities.User) error {
	_, err := r.DB.Exec(`
		INSERT OR REPLACE INTO users (id, username, email, password_hash, full_name, bio, avatar_url, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, user.ID, user.Username, user.Email, user.PasswordHash, user.FullName, user.Bio, user.AvatarURL, user.Role, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *SQLiteUserRepository) FindByID(id string) (*entities.User, error) {
	user := &entities.User{}
	err := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, full_name, bio, avatar_url, role, created_at, updated_at
		FROM users WHERE id = ?
	`, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.Bio, &user.AvatarURL, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *SQLiteUserRepository) FindByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	err := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, full_name, bio, avatar_url, role, created_at, updated_at
		FROM users WHERE email = ?
	`, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.Bio, &user.AvatarURL, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *SQLiteUserRepository) FindByUsername(username string) (*entities.User, error) {
	user := &entities.User{}
	err := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, full_name, bio, avatar_url, role, created_at, updated_at
		FROM users WHERE username = ?
	`, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.Bio, &user.AvatarURL, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *SQLiteUserRepository) ExistsByEmail(email string) (bool, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SQLiteUserRepository) ExistsByUsername(username string) (bool, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SQLiteUserRepository) GetAll() ([]*entities.User, error) {
	rows, err := r.DB.Query(`
		SELECT id, username, email, password_hash, full_name, bio, avatar_url, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.FullName, &user.Bio, &user.AvatarURL, &user.Role, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *SQLiteUserRepository) Delete(id string) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
