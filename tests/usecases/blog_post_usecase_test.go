package usecases_test

import (
	"gocleanarchitecture/entities"
	"gocleanarchitecture/frameworks/logger"
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

func (m *MockLogger) Debug(msg string, fields ...logger.LogField) {}
func (m *MockLogger) Info(msg string, fields ...logger.LogField)  {}
func (m *MockLogger) Warn(msg string, fields ...logger.LogField)  {}
func (m *MockLogger) Error(msg string, fields ...logger.LogField) {}

func TestCreateBlogPost(t *testing.T) {
	repo := &MockBlogPostRepository{blogPosts: make(map[string]*entities.BlogPost)}
	mockLogger := &MockLogger{}
	usecase := usecases.BlogPostUseCase{Repo: repo, Logger: mockLogger}

	blogPost := &entities.BlogPost{ID: "1", Title: "Test Title", Content: "Test Content"}
	err := usecase.CreateBlogPost(blogPost)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.blogPosts) != 1 {
		t.Fatalf("expected 1 blog post, got %d", len(repo.blogPosts))
	}
}

func TestGetAllBlogPosts(t *testing.T) {
	repo := &MockBlogPostRepository{blogPosts: make(map[string]*entities.BlogPost)}
	mockLogger := &MockLogger{}
	usecase := usecases.BlogPostUseCase{Repo: repo, Logger: mockLogger}

	blogPost1 := &entities.BlogPost{ID: "1", Title: "Title 1", Content: "Content 1"}
	blogPost2 := &entities.BlogPost{ID: "2", Title: "Title 2", Content: "Content 2"}

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

	blogPost := &entities.BlogPost{ID: "1", Title: "Test Title", Content: "Test Content"}
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

	blogPost := &entities.BlogPost{ID: "1", Title: "Original Title", Content: "Original Content"}
	repo.Save(blogPost)

	updatedBlogPost := &entities.BlogPost{ID: "1", Title: "Updated Title", Content: "Updated Content"}
	err := usecase.UpdateBlogPost(updatedBlogPost)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result, _ := repo.FindByID("1")
	if result.Title != "Updated Title" || result.Content != "Updated Content" {
		t.Fatalf("blog post was not updated correctly: %v", result)
	}
}

func TestDeleteBlogPost(t *testing.T) {
	repo := &MockBlogPostRepository{blogPosts: make(map[string]*entities.BlogPost)}
	mockLogger := &MockLogger{}
	usecase := usecases.BlogPostUseCase{Repo: repo, Logger: mockLogger}

	blogPost := &entities.BlogPost{ID: "1", Title: "Test Title", Content: "Test Content"}
	repo.Save(blogPost)

	err := usecase.DeleteBlogPost("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result, _ := repo.FindByID("1")
	if result != nil {
		t.Fatalf("expected blog post to be deleted, but it still exists")
	}
}
