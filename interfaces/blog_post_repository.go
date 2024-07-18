package interfaces

import "gocleanarchitecture/entities"

type BlogPostRepository interface {
	Save(blogPost entities.BlogPost) error
	FindAll() ([]entities.BlogPost, error)
}
