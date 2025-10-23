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

type SupabaseCommentRepository struct {
	URL    string
	APIKey string
	client *http.Client
}

type supabaseComment struct {
	ID         string    `json:"id"`
	BlogPostID string    `json:"blog_post_id"`
	AuthorID   string    `json:"author_id"`
	Content    string    `json:"content"`
	ParentID   string    `json:"parent_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewSupabaseCommentRepository(url, apiKey string) interfaces.CommentRepository {
	return &SupabaseCommentRepository{
		URL:    url,
		APIKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *SupabaseCommentRepository) Save(comment *entities.Comment) error {
	supabaseComment := supabaseComment{
		ID:         comment.ID,
		BlogPostID: comment.BlogPostID,
		AuthorID:   comment.AuthorID,
		Content:    comment.Content,
		ParentID:   comment.ParentID,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
	}

	jsonData, err := json.Marshal(supabaseComment)
	if err != nil {
		return fmt.Errorf("failed to marshal comment: %w", err)
	}

	req, err := http.NewRequest("POST", r.URL+"/rest/v1/comments", bytes.NewBuffer(jsonData))
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

func (r *SupabaseCommentRepository) FindByID(id string) (*entities.Comment, error) {
	url := fmt.Sprintf("%s/rest/v1/comments?id=eq.%s", r.URL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	var supabaseComments []supabaseComment
	if err := json.NewDecoder(resp.Body).Decode(&supabaseComments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(supabaseComments) == 0 {
		return nil, nil
	}

	return r.toEntity(&supabaseComments[0]), nil
}

func (r *SupabaseCommentRepository) FindByBlogPostID(blogPostID string) ([]*entities.Comment, error) {
	url := fmt.Sprintf("%s/rest/v1/comments?blog_post_id=eq.%s&order=created_at.asc", r.URL, blogPostID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	var supabaseComments []supabaseComment
	if err := json.NewDecoder(resp.Body).Decode(&supabaseComments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	comments := make([]*entities.Comment, len(supabaseComments))
	for i, sc := range supabaseComments {
		comments[i] = r.toEntity(&sc)
	}

	return comments, nil
}

func (r *SupabaseCommentRepository) FindRepliesByParentID(parentID string) ([]*entities.Comment, error) {
	url := fmt.Sprintf("%s/rest/v1/comments?parent_id=eq.%s&order=created_at.asc", r.URL, parentID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	var supabaseComments []supabaseComment
	if err := json.NewDecoder(resp.Body).Decode(&supabaseComments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	comments := make([]*entities.Comment, len(supabaseComments))
	for i, sc := range supabaseComments {
		comments[i] = r.toEntity(&sc)
	}

	return comments, nil
}

func (r *SupabaseCommentRepository) Delete(id string) error {
	url := fmt.Sprintf("%s/rest/v1/comments?id=eq.%s", r.URL, id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *SupabaseCommentRepository) GetAll() ([]*entities.Comment, error) {
	url := fmt.Sprintf("%s/rest/v1/comments?select=*&order=created_at.desc", r.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r.setHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase error %d: %s", resp.StatusCode, string(body))
	}

	var supabaseComments []supabaseComment
	if err := json.NewDecoder(resp.Body).Decode(&supabaseComments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	comments := make([]*entities.Comment, len(supabaseComments))
	for i, sc := range supabaseComments {
		comments[i] = r.toEntity(&sc)
	}

	return comments, nil
}

func (r *SupabaseCommentRepository) setHeaders(req *http.Request) {
	req.Header.Set("apikey", r.APIKey)
	req.Header.Set("Authorization", "Bearer "+r.APIKey)
	req.Header.Set("Content-Type", "application/json")
}

func (r *SupabaseCommentRepository) toEntity(sc *supabaseComment) *entities.Comment {
	return &entities.Comment{
		ID:         sc.ID,
		BlogPostID: sc.BlogPostID,
		AuthorID:   sc.AuthorID,
		Content:    sc.Content,
		ParentID:   sc.ParentID,
		CreatedAt:  sc.CreatedAt,
		UpdatedAt:  sc.UpdatedAt,
	}
}
