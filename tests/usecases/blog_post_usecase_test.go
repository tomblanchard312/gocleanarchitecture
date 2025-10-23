package usecases_test

import (
	"gocleanarchitecture/entities"
	"gocleanarchitecture/usecases"
	"testing"
)

type MockBlogPostRepository struct {
	blogPosts map[string]*entities.BlogPost
}

func (m *MockBlogPostRepository) Save(blogPost *entities.BlogPost) error {
	m.blogPosts[blogPost.ID] = blogPost
	return nil
}

func (m *MockBlogPostRepository) FindAll() ([]*entities.BlogPost, error) {
	posts := make([]*entities.BlogPost, 0, len(m.blogPosts))
	for _, post := range m.blogPosts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (m *MockBlogPostRepository) FindByID(id string) (*entities.BlogPost, error) {
	return m.blogPosts[id], nil
}

func (m *MockBlogPostRepository) Delete(id string) error {
	delete(m.blogPosts, id)
	return nil
}

type MockLogger struct{}

func (m *MockLogger) Error(msg string, fields ...interface{}) {}

func TestCreateBlogPost(t *testing.T) {
	repo := &MockBlogPostRepository{blogPosts: make(map[string]*entities.BlogPost)}
	mockLogger := &MockLogger{}
	usecase := usecases.BlogPostUseCase{Repo: repo, Logger: mockLogger}

	authorID := "user-123"
	blogPost, err := usecase.CreateBlogPost("1", "Test Title", "Test Content", authorID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.blogPosts) != 1 {
		t.Fatalf("expected 1 blog post, got %d", len(repo.blogPosts))
	}

	if blogPost.ID != "1" || blogPost.Title != "Test Title" || blogPost.Content != "Test Content" {
		t.Fatalf("blog post not created correctly: %v", blogPost)
	}

	if blogPost.AuthorID != authorID {
		t.Errorf("expected author ID to be %s, got %s", authorID, blogPost.AuthorID)
	}
}

func TestGetAllBlogPosts(t *testing.T) {
	repo := &MockBlogPostRepository{blogPosts: make(map[string]*entities.BlogPost)}
	mockLogger := &MockLogger{}
	usecase := usecases.BlogPostUseCase{Repo: repo, Logger: mockLogger}

	blogPost1 := &entities.BlogPost{ID: "1", Title: "Title 1", Content: "Content 1", AuthorID: "user-123"}
	blogPost2 := &entities.BlogPost{ID: "2", Title: "Title 2", Content: "Content 2", AuthorID: "user-456"}

	repo.Save(blogPost1)
	repo.Save(blogPost2)

	blogPosts, err := usecase.GetAllBlogPosts()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(blogPosts) != 2 {
		t.Fatalf("expected 2 blog posts, got %d", len(blogPosts))
	}
}

func TestGetBlogPost(t *testing.T) {
	repo := &MockBlogPostRepository{blogPosts: make(map[string]*entities.BlogPost)}
	mockLogger := &MockLogger{}
	usecase := usecases.BlogPostUseCase{Repo: repo, Logger: mockLogger}

	blogPost := &entities.BlogPost{ID: "1", Title: "Test Title", Content: "Test Content", AuthorID: "user-123"}
	repo.Save(blogPost)

	result, err := usecase.GetBlogPost("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "1" || result.Title != "Test Title" || result.Content != "Test Content" {
		t.Fatalf("got unexpected blog post: %v", result)
	}
}

func TestUpdateBlogPost(t *testing.T) {
	repo := &MockBlogPostRepository{blogPosts: make(map[string]*entities.BlogPost)}
	mockLogger := &MockLogger{}
	usecase := usecases.BlogPostUseCase{Repo: repo, Logger: mockLogger}

	authorID := "user-123"
	// Create initial blog post
	_, err := usecase.CreateBlogPost("1", "Original Title", "Original Content", authorID)
	if err != nil {
		t.Fatalf("expected no error creating blog post, got %v", err)
	}

	// Update the blog post as the author
	updatedBlogPost, err := usecase.UpdateBlogPost("1", "Updated Title", "Updated Content", authorID)
	if err != nil {
		t.Fatalf("expected no error updating, got %v", err)
	}

	if updatedBlogPost.Title != "Updated Title" || updatedBlogPost.Content != "Updated Content" {
		t.Fatalf("blog post was not updated correctly: %v", updatedBlogPost)
	}
}

func TestUpdateBlogPostUnauthorized(t *testing.T) {
	repo := &MockBlogPostRepository{blogPosts: make(map[string]*entities.BlogPost)}
	mockLogger := &MockLogger{}
	usecase := usecases.BlogPostUseCase{Repo: repo, Logger: mockLogger}

	authorID := "user-123"
	differentUserID := "user-456"

	// Create blog post as user-123
	_, err := usecase.CreateBlogPost("1", "Original Title", "Original Content", authorID)
	if err != nil {
		t.Fatalf("expected no error creating blog post, got %v", err)
	}

	// Try to update as different user (should fail)
	_, err = usecase.UpdateBlogPost("1", "Updated Title", "Updated Content", differentUserID)
	if err == nil {
		t.Fatal("expected error when non-author tries to update, got nil")
	}

	if err.Error() != "unauthorized: you can only update your own blog posts" {
		t.Errorf("expected authorization error, got: %v", err)
	}
}

func TestDeleteBlogPost(t *testing.T) {
	repo := &MockBlogPostRepository{blogPosts: make(map[string]*entities.BlogPost)}
	mockLogger := &MockLogger{}
	usecase := usecases.BlogPostUseCase{Repo: repo, Logger: mockLogger}

	authorID := "user-123"
	blogPost := &entities.BlogPost{ID: "1", Title: "Test Title", Content: "Test Content", AuthorID: authorID}
	repo.Save(blogPost)

	err := usecase.DeleteBlogPost("1", authorID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result, _ := repo.FindByID("1")
	if result != nil {
		t.Fatalf("expected blog post to be deleted, but it still exists")
	}
}

func TestDeleteBlogPostUnauthorized(t *testing.T) {
	repo := &MockBlogPostRepository{blogPosts: make(map[string]*entities.BlogPost)}
	mockLogger := &MockLogger{}
	usecase := usecases.BlogPostUseCase{Repo: repo, Logger: mockLogger}

	authorID := "user-123"
	differentUserID := "user-456"
	blogPost := &entities.BlogPost{ID: "1", Title: "Test Title", Content: "Test Content", AuthorID: authorID}
	repo.Save(blogPost)

	// Try to delete as different user (should fail)
	err := usecase.DeleteBlogPost("1", differentUserID)
	if err == nil {
		t.Fatal("expected error when non-author tries to delete, got nil")
	}

	if err.Error() != "unauthorized: you can only delete your own blog posts" {
		t.Errorf("expected authorization error, got: %v", err)
	}

	// Verify post still exists
	result, _ := repo.FindByID("1")
	if result == nil {
		t.Fatal("expected blog post to still exist after unauthorized delete attempt")
	}
}
