package entities_test

import (
	"gocleanarchitecture/entities"
	"testing"
	"time"
)

func TestBlogPost(t *testing.T) {
	now := time.Now()
	bp := &entities.BlogPost{
		ID:        "1",
		Title:     "Test Title",
		Content:   "Test Content",
		AuthorID:  "user-123",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if bp.ID != "1" {
		t.Errorf("Expected ID to be '1', got '%s'", bp.ID)
	}

	if bp.Title != "Test Title" {
		t.Errorf("Expected Title to be 'Test Title', got '%s'", bp.Title)
	}

	if bp.Content != "Test Content" {
		t.Errorf("Expected Content to be 'Test Content', got '%s'", bp.Content)
	}

	if bp.AuthorID != "user-123" {
		t.Errorf("Expected AuthorID to be 'user-123', got '%s'", bp.AuthorID)
	}

	if !bp.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt to be '%v', got '%v'", now, bp.CreatedAt)
	}

	if !bp.UpdatedAt.Equal(now) {
		t.Errorf("Expected UpdatedAt to be '%v', got '%v'", now, bp.UpdatedAt)
	}
}

func TestNewBlogPost(t *testing.T) {
	bp, err := entities.NewBlogPost("1", "Test Title", "Test Content", "user-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bp.ID != "1" {
		t.Errorf("Expected ID to be '1', got '%s'", bp.ID)
	}

	if bp.Title != "Test Title" {
		t.Errorf("Expected Title to be 'Test Title', got '%s'", bp.Title)
	}

	if bp.Content != "Test Content" {
		t.Errorf("Expected Content to be 'Test Content', got '%s'", bp.Content)
	}

	if bp.AuthorID != "user-123" {
		t.Errorf("Expected AuthorID to be 'user-123', got '%s'", bp.AuthorID)
	}

	if bp.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if bp.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestNewBlogPostWithoutAuthor(t *testing.T) {
	_, err := entities.NewBlogPost("1", "Test Title", "Test Content", "")
	if err == nil {
		t.Fatal("Expected error for empty author ID, got nil")
	}

	if err.Error() != "author ID cannot be empty" {
		t.Errorf("Expected 'author ID cannot be empty' error, got %v", err)
	}
}

func TestBlogPostIsAuthor(t *testing.T) {
	bp, _ := entities.NewBlogPost("1", "Test Title", "Test Content", "user-123")

	// Test with matching user ID
	if !bp.IsAuthor("user-123") {
		t.Error("Expected IsAuthor to return true for matching user ID")
	}

	// Test with different user ID
	if bp.IsAuthor("user-456") {
		t.Error("Expected IsAuthor to return false for different user ID")
	}

	// Test with empty user ID
	if bp.IsAuthor("") {
		t.Error("Expected IsAuthor to return false for empty user ID")
	}
}
