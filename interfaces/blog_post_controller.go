package interfaces

import (
	"encoding/json"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/usecases"
	"net/http"
)

type BlogPostController struct {
	BlogPostUseCase usecases.BlogPostUseCaseInterface
}

func (c *BlogPostController) CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	var blogPost entities.BlogPost
	json.NewDecoder(r.Body).Decode(&blogPost)
	err := c.BlogPostUseCase.CreateBlogPost(blogPost)
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
