package db

import (
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
)

type InMemoryBlogPostRepository struct {
	blogPosts map[string]*entities.BlogPost
}

func NewInMemoryBlogPostRepository() interfaces.BlogPostRepository {
	return &InMemoryBlogPostRepository{
		blogPosts: make(map[string]*entities.BlogPost),
	}
}

func (r *InMemoryBlogPostRepository) Save(blogPost *entities.BlogPost) error {
	r.blogPosts[blogPost.ID] = blogPost
	return nil
}

func (r *InMemoryBlogPostRepository) FindAll() ([]*entities.BlogPost, error) {
	posts := make([]*entities.BlogPost, 0, len(r.blogPosts))
	for _, post := range r.blogPosts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *InMemoryBlogPostRepository) FindByID(id string) (*entities.BlogPost, error) {
	post, ok := r.blogPosts[id]
	if !ok {
		return nil, nil
	}
	return post, nil
}

func (r *InMemoryBlogPostRepository) Delete(id string) error {
	delete(r.blogPosts, id)
	return nil
}
