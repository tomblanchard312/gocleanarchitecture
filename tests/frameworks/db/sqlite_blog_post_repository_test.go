package db_test

import (
	"gocleanarchitecture/entities"
	"gocleanarchitecture/frameworks/db/sqlite"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestSQLiteBlogPostRepository(t *testing.T) {
	// Set up a temporary database file for testing
	tempFile, err := os.CreateTemp("", "test_sqlite_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Initialize the database
	db, err := sqlite.InitDB(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	repo := sqlite.NewSQLiteBlogPostRepository(db)

	// Test Save and FindByID
	blogPost := &entities.BlogPost{
		ID:        "1",
		Title:     "Test Title",
		Content:   "Test Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = repo.Save(blogPost)
	if err != nil {
		t.Fatalf("Failed to save blog post: %v", err)
	}

	retrievedPost, err := repo.FindByID("1")
	if err != nil {
		t.Fatalf("Failed to find blog post: %v", err)
	}

	if retrievedPost.ID != blogPost.ID || retrievedPost.Title != blogPost.Title || retrievedPost.Content != blogPost.Content {
		t.Errorf("Retrieved blog post does not match the original")
	}

	// Test FindAll
	blogPost2 := &entities.BlogPost{
		ID:        "2",
		Title:     "Test Title 2",
		Content:   "Test Content 2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = repo.Save(blogPost2)
	if err != nil {
		t.Fatalf("Failed to save second blog post: %v", err)
	}

	allPosts, err := repo.FindAll()
	if err != nil {
		t.Fatalf("Failed to find all blog posts: %v", err)
	}

	if len(allPosts) != 2 {
		t.Errorf("Expected 2 blog posts, got %d", len(allPosts))
	}

	// Test Delete
	err = repo.Delete("1")
	if err != nil {
		t.Fatalf("Failed to delete blog post: %v", err)
	}

	deletedPost, err := repo.FindByID("1")
	if err != nil {
		t.Fatalf("Unexpected error when finding deleted post: %v", err)
	}
	if deletedPost != nil {
		t.Errorf("Blog post was not deleted")
	}

	allPosts, err = repo.FindAll()
	if err != nil {
		t.Fatalf("Failed to find all blog posts after deletion: %v", err)
	}

	if len(allPosts) != 1 {
		t.Errorf("Expected 1 blog post after deletion, got %d", len(allPosts))
	}
}