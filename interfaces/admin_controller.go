package interfaces

import (
	"encoding/json"
	"gocleanarchitecture/entities"
	"net/http"

	"github.com/gorilla/mux"
)

type AdminController struct {
	UserUseCase AdminUserUseCase
}

// AdminUserUseCase defines the interface for admin user management operations
type AdminUserUseCase interface {
	GetAllUsers() ([]*entities.User, error)
	GetUserByID(userID string) (*entities.User, error)
	UpdateUserRole(userID string, newRole entities.UserRole) error
	DeleteUser(userID string) error
}

func NewAdminController(userUseCase AdminUserUseCase) *AdminController {
	return &AdminController{
		UserUseCase: userUseCase,
	}
}

// GetAllUsers returns all users (admin only)
func (c *AdminController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.UserUseCase.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Sanitize all users before returning
	sanitizedUsers := make([]*entities.User, len(users))
	for i, user := range users {
		sanitizedUsers[i] = user.Sanitize()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sanitizedUsers)
}

// GetUserDetails returns details for a specific user (admin only)
func (c *AdminController) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID is required"})
		return
	}

	user, err := c.UserUseCase.GetUserByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.Sanitize())
}

// UpdateUserRole updates a user's role (admin only)
func (c *AdminController) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID is required"})
		return
	}

	var req struct {
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate role
	if err := entities.ValidateRole(req.Role); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Update role
	if err := c.UserUseCase.UpdateUserRole(userID, entities.UserRole(req.Role)); err != nil {
		if err.Error() == "user not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User role updated successfully"})
}

// DeleteUser deletes a user (admin only)
func (c *AdminController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID is required"})
		return
	}

	// Prevent admin from deleting themselves
	currentUserID := r.Context().Value("userID").(string)
	if userID == currentUserID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot delete your own account"})
		return
	}

	if err := c.UserUseCase.DeleteUser(userID); err != nil {
		if err.Error() == "user not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
