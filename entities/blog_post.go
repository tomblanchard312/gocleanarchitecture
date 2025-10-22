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
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Domain methods and business rules
func NewBlogPost(id, title, content string) (*BlogPost, error) {
	if err := validateBlogPost(id, title, content); err != nil {
		return nil, err
	}

	now := time.Now()
	return &BlogPost{
		ID:        id,
		Title:     strings.TrimSpace(title),
		Content:   strings.TrimSpace(content),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
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
