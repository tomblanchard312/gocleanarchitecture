package interfaces_test

import (
	"bytes"
	"encoding/json"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock implementation of BlogPostUseCaseInterface
type MockBlogPostUseCase struct {
	BlogPosts []entities.BlogPost
}

func (m *MockBlogPostUseCase) CreateBlogPost(blogPost entities.BlogPost) error {
	m.BlogPosts = append(m.BlogPosts, blogPost)
	return nil
}

func (m *MockBlogPostUseCase) GetAllBlogPosts() ([]entities.BlogPost, error) {
	return m.BlogPosts, nil
}

func TestCreateBlogPostHandler(t *testing.T) {
	mockUseCase := &MockBlogPostUseCase{}
	controller := interfaces.BlogPostController{BlogPostUseCase: mockUseCase}

	blogPost := entities.BlogPost{ID: "1", Title: "Test Title", Content: "Test Content"}
	body, _ := json.Marshal(blogPost)
	req, err := http.NewRequest("POST", "/blogposts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.CreateBlogPost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if len(mockUseCase.BlogPosts) != 1 {
		t.Fatalf("expected 1 blog post, got %d", len(mockUseCase.BlogPosts))
	}
}

func TestGetAllBlogPostsHandler(t *testing.T) {
	mockUseCase := &MockBlogPostUseCase{
		BlogPosts: []entities.BlogPost{
			{ID: "1", Title: "Title 1", Content: "Content 1"},
			{ID: "2", Title: "Title 2", Content: "Content 2"},
		},
	}
	controller := interfaces.BlogPostController{BlogPostUseCase: mockUseCase}

	req, err := http.NewRequest("GET", "/blogposts", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.GetAllBlogPosts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var blogPosts []entities.BlogPost
	err = json.NewDecoder(rr.Body).Decode(&blogPosts)
	if err != nil {
		t.Fatal(err)
	}

	if len(blogPosts) != 2 {
		t.Fatalf("expected 2 blog posts, got %d", len(blogPosts))
	}
}
