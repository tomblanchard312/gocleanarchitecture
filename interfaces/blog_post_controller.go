package interfaces

import (
	"encoding/json"
	"gocleanarchitecture/entities"
	"net/http"

	"github.com/gorilla/mux"
)

type BlogPostUseCase interface {
	CreateBlogPost(blogPost *entities.BlogPost) error
	GetAllBlogPosts() ([]*entities.BlogPost, error)
	GetBlogPost(id string) (*entities.BlogPost, error)
	UpdateBlogPost(blogPost *entities.BlogPost) error
	DeleteBlogPost(id string) error
}

type BlogPostController struct {
	BlogPostUseCase BlogPostUseCase
}

func (c *BlogPostController) CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	var blogPost entities.BlogPost
	err := json.NewDecoder(r.Body).Decode(&blogPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.BlogPostUseCase.CreateBlogPost(&blogPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *BlogPostController) GetAllBlogPosts(w http.ResponseWriter, r *http.Request) {
	blogPosts, err := c.BlogPostUseCase.GetAllBlogPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	json.NewEncoder(w).Encode(blogPost)
}

func (c *BlogPostController) UpdateBlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var blogPost entities.BlogPost
	err := json.NewDecoder(r.Body).Decode(&blogPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blogPost.ID = id
	err = c.BlogPostUseCase.UpdateBlogPost(&blogPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *BlogPostController) DeleteBlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := c.BlogPostUseCase.DeleteBlogPost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
