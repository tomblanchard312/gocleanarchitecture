package interfaces

import (
	"encoding/json"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/frameworks/websocket"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CommentController struct {
	CommentUseCase CommentUseCaseInterface
	WebSocketHub   *websocket.Hub
}

type CommentUseCaseInterface interface {
	CreateComment(id, blogPostID, authorID, content, parentID string) (*entities.Comment, error)
	GetCommentsByBlogPostID(blogPostID string) ([]*entities.Comment, error)
	GetRepliesByCommentID(commentID string) ([]*entities.Comment, error)
	UpdateComment(id, content, userID string) (*entities.Comment, error)
	DeleteComment(id, userID string) error
}

func NewCommentController(commentUseCase CommentUseCaseInterface) *CommentController {
	return &CommentController{
		CommentUseCase: commentUseCase,
	}
}

// CreateComment handles POST /blogposts/{blogPostId}/comments
func (c *CommentController) CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogPostID := vars["blogPostId"]

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
		return
	}

	var req struct {
		Content  string `json:"content"`
		ParentID string `json:"parent_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Generate UUID for comment
	commentID := uuid.New().String()

	comment, err := c.CommentUseCase.CreateComment(commentID, blogPostID, userID, req.Content, req.ParentID)
	if err != nil {
		if err.Error() == "blog post not found" || err.Error() == "author not found" || err.Error() == "parent comment not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Broadcast new comment via WebSocket
	if c.WebSocketHub != nil {
		c.WebSocketHub.BroadcastJSON(websocket.MessageTypeNewComment, comment)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// GetCommentsByBlogPost handles GET /blogposts/{blogPostId}/comments
func (c *CommentController) GetCommentsByBlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogPostID := vars["blogPostId"]

	comments, err := c.CommentUseCase.GetCommentsByBlogPostID(blogPostID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// GetRepliesByComment handles GET /comments/{commentId}/replies
func (c *CommentController) GetRepliesByComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["commentId"]

	replies, err := c.CommentUseCase.GetRepliesByCommentID(commentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(replies)
}

// UpdateComment handles PUT /comments/{commentId}
func (c *CommentController) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["commentId"]

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
		return
	}

	var req struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	comment, err := c.CommentUseCase.UpdateComment(commentID, req.Content, userID)
	if err != nil {
		if err.Error() == "comment not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if err.Error() == "unauthorized: only the author can update this comment" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// DeleteComment handles DELETE /comments/{commentId}
func (c *CommentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["commentId"]

	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
		return
	}

	err := c.CommentUseCase.DeleteComment(commentID, userID)
	if err != nil {
		if err.Error() == "comment not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if err.Error() == "unauthorized: only the author or admin can delete this comment" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment deleted successfully"})
}
