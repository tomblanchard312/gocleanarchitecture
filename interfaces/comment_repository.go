package interfaces

import "gocleanarchitecture/entities"

type CommentRepository interface {
	Save(comment *entities.Comment) error
	FindByID(id string) (*entities.Comment, error)
	FindByBlogPostID(blogPostID string) ([]*entities.Comment, error)
	FindRepliesByParentID(parentID string) ([]*entities.Comment, error)
	Delete(id string) error
	GetAll() ([]*entities.Comment, error)
}
