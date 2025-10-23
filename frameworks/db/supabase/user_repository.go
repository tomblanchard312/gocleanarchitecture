package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
	"io"
	"net/http"
	"time"
)

type SupabaseUserRepository struct {
	URL    string
	APIKey string
	client *http.Client
}

type supabaseUser struct {
	ID           string            `json:"id"`
	Username     string            `json:"username"`
	Email        string            `json:"email"`
	PasswordHash string            `json:"password_hash"`
	FullName     string            `json:"full_name"`
	Bio          string            `json:"bio"`
	AvatarURL    string            `json:"avatar_url"`
	Role         entities.UserRole `json:"role"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func NewSupabaseUserRepository(url, apiKey string) interfaces.UserRepository {
	return &SupabaseUserRepository{
		URL:    url,
		APIKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *SupabaseUserRepository) Save(user *entities.User) error {
	supabaseUser := supabaseUser{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		Bio:          user.Bio,
		AvatarURL:    user.AvatarURL,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	jsonData, err := json.Marshal(supabaseUser)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	req, err := http.NewRequest("POST", r.URL+"/rest/v1/users", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)
	req.Header.Set("Prefer", "resolution=merge-duplicates")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *SupabaseUserRepository) FindByID(id string) (*entities.User, error) {
	req, err := http.NewRequest("GET", r.URL+"/rest/v1/users?id=eq."+id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	var supabaseUsers []supabaseUser
	if err := json.NewDecoder(resp.Body).Decode(&supabaseUsers); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(supabaseUsers) == 0 {
		return nil, nil
	}

	return r.toEntity(&supabaseUsers[0]), nil
}

func (r *SupabaseUserRepository) FindByEmail(email string) (*entities.User, error) {
	req, err := http.NewRequest("GET", r.URL+"/rest/v1/users?email=eq."+email, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	var supabaseUsers []supabaseUser
	if err := json.NewDecoder(resp.Body).Decode(&supabaseUsers); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(supabaseUsers) == 0 {
		return nil, nil
	}

	return r.toEntity(&supabaseUsers[0]), nil
}

func (r *SupabaseUserRepository) FindByUsername(username string) (*entities.User, error) {
	req, err := http.NewRequest("GET", r.URL+"/rest/v1/users?username=eq."+username, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	var supabaseUsers []supabaseUser
	if err := json.NewDecoder(resp.Body).Decode(&supabaseUsers); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(supabaseUsers) == 0 {
		return nil, nil
	}

	return r.toEntity(&supabaseUsers[0]), nil
}

func (r *SupabaseUserRepository) ExistsByEmail(email string) (bool, error) {
	user, err := r.FindByEmail(email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func (r *SupabaseUserRepository) ExistsByUsername(username string) (bool, error) {
	user, err := r.FindByUsername(username)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func (r *SupabaseUserRepository) GetAll() ([]*entities.User, error) {
	url := fmt.Sprintf("%s/rest/v1/users?select=*&order=created_at.desc", r.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	var supabaseUsers []supabaseUser
	if err := json.NewDecoder(resp.Body).Decode(&supabaseUsers); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	users := make([]*entities.User, len(supabaseUsers))
	for i, su := range supabaseUsers {
		users[i] = r.toEntity(&su)
	}

	return users, nil
}

func (r *SupabaseUserRepository) Delete(id string) error {
	url := fmt.Sprintf("%s/rest/v1/users?id=eq.%s", r.URL, id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *SupabaseUserRepository) setHeaders(req *http.Request) {
	req.Header.Set("apikey", r.APIKey)
	req.Header.Set("Authorization", "Bearer "+r.APIKey)
	req.Header.Set("Content-Type", "application/json")
}

func (r *SupabaseUserRepository) toEntity(su *supabaseUser) *entities.User {
	return &entities.User{
		ID:           su.ID,
		Username:     su.Username,
		Email:        su.Email,
		PasswordHash: su.PasswordHash,
		FullName:     su.FullName,
		Bio:          su.Bio,
		AvatarURL:    su.AvatarURL,
		Role:         su.Role,
		CreatedAt:    su.CreatedAt,
		UpdatedAt:    su.UpdatedAt,
	}
}
