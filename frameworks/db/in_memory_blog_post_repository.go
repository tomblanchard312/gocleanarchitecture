package db

import (
	"errors"
	"gocleanarchitecture/entities"
)

type InMemoryBlogPostRepository struct {
	blogPosts map[string]entities.BlogPost
}

func NewInMemoryBlogPostRepository() *InMemoryBlogPostRepository {
	return &InMemoryBlogPostRepository{
		blogPosts: make(map[string]entities.BlogPost),
	}
}

func (r *InMemoryBlogPostRepository) Save(blogPost entities.BlogPost) error {
	if blogPost.ID == "" {
		return errors.New("blog post ID cannot be empty")
	}
	r.blogPosts[blogPost.ID] = blogPost
	return nil
}

func (r *InMemoryBlogPostRepository) FindAll() ([]entities.BlogPost, error) {
	posts := make([]entities.BlogPost, 0, len(r.blogPosts))
	for _, post := range r.blogPosts {
		posts = append(posts, post)
	}
	return posts, nil
}
