package usecases

import (
	"errors"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
)

type CommentUseCaseInterface interface {
	CreateComment(id, blogPostID, authorID, content, parentID string) (*entities.Comment, error)
	GetCommentsByBlogPostID(blogPostID string) ([]*entities.Comment, error)
	GetRepliesByCommentID(commentID string) ([]*entities.Comment, error)
	UpdateComment(id, content, userID string) (*entities.Comment, error)
	DeleteComment(id, userID string) error
}

type CommentUseCase struct {
	CommentRepo  interfaces.CommentRepository
	BlogPostRepo interfaces.BlogPostRepository
	UserRepo     interfaces.UserRepository
	Logger       Logger
}

func NewCommentUseCase(commentRepo interfaces.CommentRepository, blogPostRepo interfaces.BlogPostRepository, userRepo interfaces.UserRepository, logger Logger) *CommentUseCase {
	return &CommentUseCase{
		CommentRepo:  commentRepo,
		BlogPostRepo: blogPostRepo,
		UserRepo:     userRepo,
		Logger:       logger,
	}
}

// CreateComment creates a new comment on a blog post
func (uc *CommentUseCase) CreateComment(id, blogPostID, authorID, content, parentID string) (*entities.Comment, error) {
	// Validate blog post exists
	blogPost, err := uc.BlogPostRepo.FindByID(blogPostID)
	if err != nil {
		uc.Logger.Error("Failed to fetch blog post", map[string]interface{}{
			"error":      err.Error(),
			"blogPostID": blogPostID,
		})
		return nil, errors.New("failed to validate blog post")
	}

	if blogPost == nil {
		return nil, errors.New("blog post not found")
	}

	// Validate author exists
	author, err := uc.UserRepo.FindByID(authorID)
	if err != nil {
		uc.Logger.Error("Failed to fetch author", map[string]interface{}{
			"error":    err.Error(),
			"authorID": authorID,
		})
		return nil, errors.New("failed to validate author")
	}

	if author == nil {
		return nil, errors.New("author not found")
	}

	// If it's a reply, validate parent comment exists
	if parentID != "" {
		parentComment, err := uc.CommentRepo.FindByID(parentID)
		if err != nil {
			uc.Logger.Error("Failed to fetch parent comment", map[string]interface{}{
				"error":    err.Error(),
				"parentID": parentID,
			})
			return nil, errors.New("failed to validate parent comment")
		}

		if parentComment == nil {
			return nil, errors.New("parent comment not found")
		}

		// Ensure parent comment is on the same blog post
		if parentComment.BlogPostID != blogPostID {
			return nil, errors.New("parent comment is not on the same blog post")
		}
	}

	// Create comment
	comment, err := entities.NewComment(id, blogPostID, authorID, content, parentID)
	if err != nil {
		return nil, err
	}

	// Save comment
	if err := uc.CommentRepo.Save(comment); err != nil {
		uc.Logger.Error("Failed to save comment", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.New("failed to create comment")
	}

	return comment, nil
}

// GetCommentsByBlogPostID retrieves all comments for a specific blog post
func (uc *CommentUseCase) GetCommentsByBlogPostID(blogPostID string) ([]*entities.Comment, error) {
	if blogPostID == "" {
		return nil, errors.New("blog post ID is required")
	}

	comments, err := uc.CommentRepo.FindByBlogPostID(blogPostID)
	if err != nil {
		uc.Logger.Error("Failed to fetch comments", map[string]interface{}{
			"error":      err.Error(),
			"blogPostID": blogPostID,
		})
		return nil, errors.New("failed to fetch comments")
	}

	return comments, nil
}

// GetRepliesByCommentID retrieves all replies to a specific comment
func (uc *CommentUseCase) GetRepliesByCommentID(commentID string) ([]*entities.Comment, error) {
	if commentID == "" {
		return nil, errors.New("comment ID is required")
	}

	replies, err := uc.CommentRepo.FindRepliesByParentID(commentID)
	if err != nil {
		uc.Logger.Error("Failed to fetch replies", map[string]interface{}{
			"error":     err.Error(),
			"commentID": commentID,
		})
		return nil, errors.New("failed to fetch replies")
	}

	return replies, nil
}

// UpdateComment updates a comment's content (only by the author)
func (uc *CommentUseCase) UpdateComment(id, content, userID string) (*entities.Comment, error) {
	if id == "" {
		return nil, errors.New("comment ID is required")
	}

	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Fetch comment
	comment, err := uc.CommentRepo.FindByID(id)
	if err != nil {
		uc.Logger.Error("Failed to fetch comment", map[string]interface{}{
			"error":     err.Error(),
			"commentID": id,
		})
		return nil, errors.New("failed to fetch comment")
	}

	if comment == nil {
		return nil, errors.New("comment not found")
	}

	// Check if user is the author
	if !comment.IsAuthor(userID) {
		return nil, errors.New("unauthorized: only the author can update this comment")
	}

	// Update content
	if err := comment.Update(content); err != nil {
		return nil, err
	}

	// Save updated comment
	if err := uc.CommentRepo.Save(comment); err != nil {
		uc.Logger.Error("Failed to update comment", map[string]interface{}{
			"error":     err.Error(),
			"commentID": id,
		})
		return nil, errors.New("failed to update comment")
	}

	return comment, nil
}

// DeleteComment deletes a comment (only by the author or admin)
func (uc *CommentUseCase) DeleteComment(id, userID string) error {
	if id == "" {
		return errors.New("comment ID is required")
	}

	if userID == "" {
		return errors.New("user ID is required")
	}

	// Fetch comment
	comment, err := uc.CommentRepo.FindByID(id)
	if err != nil {
		uc.Logger.Error("Failed to fetch comment", map[string]interface{}{
			"error":     err.Error(),
			"commentID": id,
		})
		return errors.New("failed to fetch comment")
	}

	if comment == nil {
		return errors.New("comment not found")
	}

	// Check if user is the author or admin
	user, err := uc.UserRepo.FindByID(userID)
	if err != nil {
		uc.Logger.Error("Failed to fetch user", map[string]interface{}{
			"error":  err.Error(),
			"userID": userID,
		})
		return errors.New("failed to verify permissions")
	}

	if user == nil {
		return errors.New("user not found")
	}

	if !comment.IsAuthor(userID) && !user.IsAdmin() {
		return errors.New("unauthorized: only the author or admin can delete this comment")
	}

	// Delete comment
	if err := uc.CommentRepo.Delete(id); err != nil {
		uc.Logger.Error("Failed to delete comment", map[string]interface{}{
			"error":     err.Error(),
			"commentID": id,
		})
		return errors.New("failed to delete comment")
	}

	return nil
}
