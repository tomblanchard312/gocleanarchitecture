package db

import (
	"gocleanarchitecture/entities"
	"testing"
)

func TestSaveAndFindAll(t *testing.T) {
	repo := NewInMemoryBlogPostRepository()

	blogPost1 := entities.BlogPost{ID: "1", Title: "Title 1", Content: "Content 1"}
	blogPost2 := entities.BlogPost{ID: "2", Title: "Title 2", Content: "Content 2"}

	err := repo.Save(blogPost1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = repo.Save(blogPost2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	blogPosts, err := repo.FindAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(blogPosts) != 2 {
		t.Fatalf("expected 2 blog posts, got %d", len(blogPosts))
	}
}
