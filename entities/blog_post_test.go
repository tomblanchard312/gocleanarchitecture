package entities

import "testing"

func TestBlogPostCreation(t *testing.T) {
	blogPost := BlogPost{ID: "1", Title: "Test Title", Content: "Test Content"}

	if blogPost.ID != "1" {
		t.Fatalf("expected ID to be '1', got %s", blogPost.ID)
	}

	if blogPost.Title != "Test Title" {
		t.Fatalf("expected Title to be 'Test Title', got %s", blogPost.Title)
	}

	if blogPost.Content != "Test Content" {
		t.Fatalf("expected Content to be 'Test Content', got %s", blogPost.Content)
	}
}
