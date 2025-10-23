package db

import (
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
	"sync"
)

type InMemoryCommentRepository struct {
	comments map[string]*entities.Comment
	mu       sync.RWMutex
}

func NewInMemoryCommentRepository() interfaces.CommentRepository {
	return &InMemoryCommentRepository{
		comments: make(map[string]*entities.Comment),
	}
}

func (r *InMemoryCommentRepository) Save(comment *entities.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create a copy to avoid external modifications
	commentCopy := *comment
	r.comments[comment.ID] = &commentCopy
	return nil
}

func (r *InMemoryCommentRepository) FindByID(id string) (*entities.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[id]
	if !exists {
		return nil, nil
	}

	// Return a copy to avoid external modifications
	commentCopy := *comment
	return &commentCopy, nil
}

func (r *InMemoryCommentRepository) FindByBlogPostID(blogPostID string) ([]*entities.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var comments []*entities.Comment
	for _, comment := range r.comments {
		if comment.BlogPostID == blogPostID {
			// Return a copy to avoid external modifications
			commentCopy := *comment
			comments = append(comments, &commentCopy)
		}
	}

	return comments, nil
}

func (r *InMemoryCommentRepository) FindRepliesByParentID(parentID string) ([]*entities.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var comments []*entities.Comment
	for _, comment := range r.comments {
		if comment.ParentID == parentID {
			// Return a copy to avoid external modifications
			commentCopy := *comment
			comments = append(comments, &commentCopy)
		}
	}

	return comments, nil
}

func (r *InMemoryCommentRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.comments, id)
	return nil
}

func (r *InMemoryCommentRepository) GetAll() ([]*entities.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comments := make([]*entities.Comment, 0, len(r.comments))
	for _, comment := range r.comments {
		// Return a copy to avoid external modifications
		commentCopy := *comment
		comments = append(comments, &commentCopy)
	}

	return comments, nil
}
