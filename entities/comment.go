package entities

import (
	"errors"
	"strings"
	"time"
)

type Comment struct {
	ID         string
	BlogPostID string
	AuthorID   string
	Content    string
	ParentID   string // For nested comments/replies (empty string if top-level)
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewComment creates a new comment with validation
func NewComment(id, blogPostID, authorID, content, parentID string) (*Comment, error) {
	if err := validateCommentData(blogPostID, authorID, content); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Comment{
		ID:         id,
		BlogPostID: strings.TrimSpace(blogPostID),
		AuthorID:   strings.TrimSpace(authorID),
		Content:    strings.TrimSpace(content),
		ParentID:   strings.TrimSpace(parentID),
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// Update updates a comment's content
func (c *Comment) Update(content string) error {
	content = strings.TrimSpace(content)

	if content == "" {
		return errors.New("comment content cannot be empty")
	}

	if len(content) > 1000 {
		return errors.New("comment content cannot exceed 1000 characters")
	}

	c.Content = content
	c.UpdatedAt = time.Now()
	return nil
}

// IsAuthor checks if a given user ID is the author of the comment
func (c *Comment) IsAuthor(userID string) bool {
	return c.AuthorID == userID
}

// IsReply checks if the comment is a reply to another comment
func (c *Comment) IsReply() bool {
	return c.ParentID != ""
}

// Private validation functions

func validateCommentData(blogPostID, authorID, content string) error {
	if blogPostID == "" {
		return errors.New("blog post ID is required")
	}

	if authorID == "" {
		return errors.New("author ID is required")
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("comment content cannot be empty")
	}

	if len(content) < 1 {
		return errors.New("comment content must be at least 1 character")
	}

	if len(content) > 1000 {
		return errors.New("comment content cannot exceed 1000 characters")
	}

	return nil
}
