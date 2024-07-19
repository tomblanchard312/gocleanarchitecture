package interfaces_test

import (
	"bytes"
	"encoding/json"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type MockBlogPostUseCase struct {
	blogPosts map[string]*entities.BlogPost
}

func (m *MockBlogPostUseCase) CreateBlogPost(blogPost *entities.BlogPost) error {
	m.blogPosts[blogPost.ID] = blogPost
	return nil
}

func (m *MockBlogPostUseCase) GetAllBlogPosts() ([]*entities.BlogPost, error) {
	posts := make([]*entities.BlogPost, 0, len(m.blogPosts))
	for _, post := range m.blogPosts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (m *MockBlogPostUseCase) GetBlogPost(id string) (*entities.BlogPost, error) {
	return m.blogPosts[id], nil
}

func (m *MockBlogPostUseCase) UpdateBlogPost(blogPost *entities.BlogPost) error {
	m.blogPosts[blogPost.ID] = blogPost
	return nil
}

func (m *MockBlogPostUseCase) DeleteBlogPost(id string) error {
	delete(m.blogPosts, id)
	return nil
}

func TestCreateBlogPostHandler(t *testing.T) {
	mockUseCase := &MockBlogPostUseCase{blogPosts: make(map[string]*entities.BlogPost)}
	controller := interfaces.BlogPostController{BlogPostUseCase: mockUseCase}

	blogPost := &entities.BlogPost{ID: "1", Title: "Test Title", Content: "Test Content"}
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

	if len(mockUseCase.blogPosts) != 1 {
		t.Errorf("expected 1 blog post, got %d", len(mockUseCase.blogPosts))
	}
}

func TestGetAllBlogPostsHandler(t *testing.T) {
	mockUseCase := &MockBlogPostUseCase{
		blogPosts: map[string]*entities.BlogPost{
			"1": {ID: "1", Title: "Title 1", Content: "Content 1"},
			"2": {ID: "2", Title: "Title 2", Content: "Content 2"},
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

	var blogPosts []*entities.BlogPost
	err = json.NewDecoder(rr.Body).Decode(&blogPosts)
	if err != nil {
		t.Fatal(err)
	}

	if len(blogPosts) != 2 {
		t.Errorf("expected 2 blog posts, got %d", len(blogPosts))
	}
}

func TestGetBlogPostHandler(t *testing.T) {
	mockUseCase := &MockBlogPostUseCase{
		blogPosts: map[string]*entities.BlogPost{
			"1": {ID: "1", Title: "Test Title", Content: "Test Content"},
		},
	}
	controller := interfaces.BlogPostController{BlogPostUseCase: mockUseCase}

	req, err := http.NewRequest("GET", "/blogposts/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.GetBlogPost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var blogPost entities.BlogPost
	err = json.NewDecoder(rr.Body).Decode(&blogPost)
	if err != nil {
		t.Fatal(err)
	}

	if blogPost.ID != "1" || blogPost.Title != "Test Title" || blogPost.Content != "Test Content" {
		t.Errorf("handler returned unexpected body: got %v", blogPost)
	}
}

func TestUpdateBlogPostHandler(t *testing.T) {
	mockUseCase := &MockBlogPostUseCase{
		blogPosts: map[string]*entities.BlogPost{
			"1": {ID: "1", Title: "Original Title", Content: "Original Content"},
		},
	}
	controller := interfaces.BlogPostController{BlogPostUseCase: mockUseCase}

	updatedBlogPost := &entities.BlogPost{ID: "1", Title: "Updated Title", Content: "Updated Content"}
	body, _ := json.Marshal(updatedBlogPost)
	req, err := http.NewRequest("PUT", "/blogposts/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.UpdateBlogPost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	updatedPost := mockUseCase.blogPosts["1"]
	if updatedPost.Title != "Updated Title" || updatedPost.Content != "Updated Content" {
		t.Errorf("blog post was not updated correctly: got %v", updatedPost)
	}
}

func TestDeleteBlogPostHandler(t *testing.T) {
	mockUseCase := &MockBlogPostUseCase{
		blogPosts: map[string]*entities.BlogPost{
			"1": {ID: "1", Title: "Test Title", Content: "Test Content"},
		},
	}
	controller := interfaces.BlogPostController{BlogPostUseCase: mockUseCase}

	req, err := http.NewRequest("DELETE", "/blogposts/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.DeleteBlogPost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	if _, exists := mockUseCase.blogPosts["1"]; exists {
		t.Errorf("blog post was not deleted")
	}
}
