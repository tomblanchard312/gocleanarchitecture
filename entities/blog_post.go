package entities

import (
	"errors"
	"strings"
	"time"
)

type BlogPost struct {
	ID        string
	Title     string
	Content   string
	AuthorID  string // User ID of the author
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Domain methods and business rules
func NewBlogPost(id, title, content, authorID string) (*BlogPost, error) {
	if err := validateBlogPost(id, title, content); err != nil {
		return nil, err
	}

	if strings.TrimSpace(authorID) == "" {
		return nil, errors.New("author ID cannot be empty")
	}

	now := time.Now()
	return &BlogPost{
		ID:        id,
		Title:     strings.TrimSpace(title),
		Content:   strings.TrimSpace(content),
		AuthorID:  authorID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// IsAuthor checks if the given user ID is the author of this blog post
func (bp *BlogPost) IsAuthor(userID string) bool {
	return bp.AuthorID == userID
}

func (bp *BlogPost) Update(title, content string) error {
	if err := validateBlogPost(bp.ID, title, content); err != nil {
		return err
	}

	bp.Title = strings.TrimSpace(title)
	bp.Content = strings.TrimSpace(content)
	bp.UpdatedAt = time.Now()
	return nil
}

func (bp *BlogPost) IsValid() bool {
	return validateBlogPost(bp.ID, bp.Title, bp.Content) == nil
}

func validateBlogPost(id, title, content string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("blog post ID cannot be empty")
	}
	if strings.TrimSpace(title) == "" {
		return errors.New("blog post title cannot be empty")
	}
	if strings.TrimSpace(content) == "" {
		return errors.New("blog post content cannot be empty")
	}
	if len(title) > 200 {
		return errors.New("blog post title cannot exceed 200 characters")
	}
	return nil
}
