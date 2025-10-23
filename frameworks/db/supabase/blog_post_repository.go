package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
	"io"
	"net/http"
	"time"
)

type SupabaseBlogPostRepository struct {
	URL    string
	APIKey string
	client *http.Client
}

type supabaseBlogPost struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewSupabaseBlogPostRepository(url, apiKey string) interfaces.BlogPostRepository {
	return &SupabaseBlogPostRepository{
		URL:    url,
		APIKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *SupabaseBlogPostRepository) Save(blogPost *entities.BlogPost) error {
	supabasePost := supabaseBlogPost{
		ID:        blogPost.ID,
		Title:     blogPost.Title,
		Content:   blogPost.Content,
		AuthorID:  blogPost.AuthorID,
		CreatedAt: blogPost.CreatedAt,
		UpdatedAt: blogPost.UpdatedAt,
	}

	jsonData, err := json.Marshal(supabasePost)
	if err != nil {
		return fmt.Errorf("failed to marshal blog post: %w", err)
	}

	req, err := http.NewRequest("POST", r.URL+"/rest/v1/blog_posts", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)
	req.Header.Set("Prefer", "resolution=merge-duplicates")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *SupabaseBlogPostRepository) FindAll() ([]*entities.BlogPost, error) {
	req, err := http.NewRequest("GET", r.URL+"/rest/v1/blog_posts", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	var supabasePosts []supabaseBlogPost
	if err := json.NewDecoder(resp.Body).Decode(&supabasePosts); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	blogPosts := make([]*entities.BlogPost, len(supabasePosts))
	for i, sp := range supabasePosts {
		blogPosts[i] = &entities.BlogPost{
			ID:        sp.ID,
			Title:     sp.Title,
			Content:   sp.Content,
			AuthorID:  sp.AuthorID,
			CreatedAt: sp.CreatedAt,
			UpdatedAt: sp.UpdatedAt,
		}
	}

	return blogPosts, nil
}

func (r *SupabaseBlogPostRepository) FindByID(id string) (*entities.BlogPost, error) {
	req, err := http.NewRequest("GET", r.URL+"/rest/v1/blog_posts?id=eq."+id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	var supabasePosts []supabaseBlogPost
	if err := json.NewDecoder(resp.Body).Decode(&supabasePosts); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(supabasePosts) == 0 {
		return nil, nil // Not found
	}

	sp := supabasePosts[0]
	return &entities.BlogPost{
		ID:        sp.ID,
		Title:     sp.Title,
		Content:   sp.Content,
		AuthorID:  sp.AuthorID,
		CreatedAt: sp.CreatedAt,
		UpdatedAt: sp.UpdatedAt,
	}, nil
}

func (r *SupabaseBlogPostRepository) Delete(id string) error {
	req, err := http.NewRequest("DELETE", r.URL+"/rest/v1/blog_posts?id=eq."+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *SupabaseBlogPostRepository) setHeaders(req *http.Request) {
	req.Header.Set("apikey", r.APIKey)
	req.Header.Set("Authorization", "Bearer "+r.APIKey)
	req.Header.Set("Content-Type", "application/json")
}
