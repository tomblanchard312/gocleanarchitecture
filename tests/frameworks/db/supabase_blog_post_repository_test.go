package db_test

import (
	"gocleanarchitecture/frameworks/db/supabase"
	"testing"
)

// This is an integration test that requires actual Supabase credentials
// For unit testing, we use the in-memory repository in other tests
func TestSupabaseBlogPostRepository_Integration(t *testing.T) {
	// Skip if not running integration tests
	t.Skip("Skipping Supabase integration test - requires live database")

	// Example usage (uncomment and set real credentials for integration testing):
	/*
		url := "https://your-project.supabase.co"
		apiKey := "your-anon-key"

		repo := supabase.NewSupabaseBlogPostRepository(url, apiKey)

		// Test create
		blogPost, err := entities.NewBlogPost("test-1", "Test Title", "Test Content")
		if err != nil {
			t.Fatalf("Failed to create blog post: %v", err)
		}

		err = repo.Save(blogPost)
		if err != nil {
			t.Fatalf("Failed to save blog post: %v", err)
		}

		// Test retrieve
		retrieved, err := repo.FindByID("test-1")
		if err != nil {
			t.Fatalf("Failed to find blog post: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Blog post not found")
		}

		if retrieved.Title != "Test Title" {
			t.Errorf("Expected title 'Test Title', got %s", retrieved.Title)
		}

		// Test delete
		err = repo.Delete("test-1")
		if err != nil {
			t.Fatalf("Failed to delete blog post: %v", err)
		}
	*/
}

func TestSupabaseBlogPostRepository_Creation(t *testing.T) {
	// Test that we can create a repository instance
	repo := supabase.NewSupabaseBlogPostRepository("https://example.supabase.co", "test-key")
	if repo == nil {
		t.Fatal("Expected repository instance, got nil")
	}
}
