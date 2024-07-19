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

	if !bp.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt to be '%v', got '%v'", now, bp.CreatedAt)
	}

	if !bp.UpdatedAt.Equal(now) {
		t.Errorf("Expected UpdatedAt to be '%v', got '%v'", now, bp.UpdatedAt)
	}
}
