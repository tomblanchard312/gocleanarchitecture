package usecases

import (
	"gocleanarchitecture/entities"
	"gocleanarchitecture/frameworks/logger"
	"testing"
)

type MockBlogPostRepository struct {
	blogPosts []entities.BlogPost
}

func (m *MockBlogPostRepository) Save(blogPost entities.BlogPost) error {
	m.blogPosts = append(m.blogPosts, blogPost)
	return nil
}

func (m *MockBlogPostRepository) FindAll() ([]entities.BlogPost, error) {
	return m.blogPosts, nil
}

type MockLogger struct{}

func (m *MockLogger) Debug(msg string, fields ...logger.LogField) {}
func (m *MockLogger) Info(msg string, fields ...logger.LogField)  {}
func (m *MockLogger) Warn(msg string, fields ...logger.LogField)  {}
func (m *MockLogger) Error(msg string, fields ...logger.LogField) {}

func TestCreateBlogPost(t *testing.T) {
	repo := &MockBlogPostRepository{}
	mockLogger := &MockLogger{}
	usecase := BlogPostUseCase{Repo: repo, Logger: mockLogger}

	blogPost := entities.BlogPost{ID: "1", Title: "Test Title", Content: "Test Content"}
	err := usecase.CreateBlogPost(blogPost)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.blogPosts) != 1 {
		t.Fatalf("expected 1 blog post, got %d", len(repo.blogPosts))
	}
}

func TestGetAllBlogPosts(t *testing.T) {
	repo := &MockBlogPostRepository{}
	mockLogger := &MockLogger{}
	usecase := BlogPostUseCase{Repo: repo, Logger: mockLogger}

	blogPost1 := entities.BlogPost{ID: "1", Title: "Title 1", Content: "Content 1"}
	blogPost2 := entities.BlogPost{ID: "2", Title: "Title 2", Content: "Content 2"}

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
