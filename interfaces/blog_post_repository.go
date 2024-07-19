package interfaces

import "gocleanarchitecture/entities"

type BlogPostRepository interface {
	Save(blogPost *entities.BlogPost) error
	FindAll() ([]*entities.BlogPost, error)
	FindByID(id string) (*entities.BlogPost, error)
	Delete(id string) error
}
