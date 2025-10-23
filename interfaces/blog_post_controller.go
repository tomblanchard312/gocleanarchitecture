package interfaces

import (
	"encoding/json"
	"gocleanarchitecture/entities"
	"net/http"

	"github.com/gorilla/mux"
)

type BlogPostUseCase interface {
	CreateBlogPost(id, title, content, authorID string) (*entities.BlogPost, error)
	GetAllBlogPosts() ([]*entities.BlogPost, error)
	GetBlogPost(id string) (*entities.BlogPost, error)
	UpdateBlogPost(id, title, content, userID string) (*entities.BlogPost, error)
	DeleteBlogPost(id, userID string) error
}

type BlogPostController struct {
	BlogPostUseCase BlogPostUseCase
}

func (c *BlogPostController) CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	blogPost, err := c.BlogPostUseCase.CreateBlogPost(request.ID, request.Title, request.Content, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blogPost)
}

func (c *BlogPostController) GetAllBlogPosts(w http.ResponseWriter, r *http.Request) {
	blogPosts, err := c.BlogPostUseCase.GetAllBlogPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogPosts)
}

func (c *BlogPostController) GetBlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	blogPost, err := c.BlogPostUseCase.GetBlogPost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if blogPost == nil {
		http.Error(w, "Blog post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogPost)
}

func (c *BlogPostController) UpdateBlogPost(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var request struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	blogPost, err := c.BlogPostUseCase.UpdateBlogPost(id, request.Title, request.Content, userID)
	if err != nil {
		// Check for authorization error
		if err.Error() == "unauthorized: you can only update your own blog posts" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(blogPost)
}

func (c *BlogPostController) DeleteBlogPost(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	err := c.BlogPostUseCase.DeleteBlogPost(id, userID)
	if err != nil {
		// Check for authorization error
		if err.Error() == "unauthorized: you can only delete your own blog posts" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
