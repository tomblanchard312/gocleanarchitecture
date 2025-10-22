package usecases

import (
	"errors"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
)

// Domain service interface (abstraction for cross-cutting concerns)
type Logger interface {
	Error(msg string, fields ...interface{})
}

type BlogPostUseCaseInterface interface {
	CreateBlogPost(id, title, content string) (*entities.BlogPost, error)
	GetAllBlogPosts() ([]*entities.BlogPost, error)
	GetBlogPost(id string) (*entities.BlogPost, error)
	UpdateBlogPost(id, title, content string) (*entities.BlogPost, error)
	DeleteBlogPost(id string) error
}

type BlogPostUseCase struct {
	Repo   interfaces.BlogPostRepository
	Logger Logger
}

func NewBlogPostUseCase(repo interfaces.BlogPostRepository, logger Logger) BlogPostUseCaseInterface {
	return &BlogPostUseCase{
		Repo:   repo,
		Logger: logger,
	}
}

func (u *BlogPostUseCase) CreateBlogPost(id, title, content string) (*entities.BlogPost, error) {
	// Use domain factory method
	blogPost, err := entities.NewBlogPost(id, title, content)
	if err != nil {
		return nil, err
	}

	err = u.Repo.Save(blogPost)
	if err != nil {
		u.Logger.Error("Failed to create blog post", "error", err)
		return nil, err
	}
	return blogPost, nil
}

func (u *BlogPostUseCase) GetAllBlogPosts() ([]*entities.BlogPost, error) {
	blogPosts, err := u.Repo.FindAll()
	if err != nil {
		u.Logger.Error("Failed to get all blog posts", "error", err)
		return nil, err
	}
	return blogPosts, nil
}

func (u *BlogPostUseCase) GetBlogPost(id string) (*entities.BlogPost, error) {
	blogPost, err := u.Repo.FindByID(id)
	if err != nil {
		u.Logger.Error("Failed to get blog post", "error", err, "id", id)
		return nil, err
	}
	return blogPost, nil
}

func (u *BlogPostUseCase) UpdateBlogPost(id, title, content string) (*entities.BlogPost, error) {
	// Get existing blog post
	blogPost, err := u.Repo.FindByID(id)
	if err != nil {
		u.Logger.Error("Failed to find blog post for update", "error", err, "id", id)
		return nil, err
	}
	if blogPost == nil {
		return nil, errors.New("blog post not found")
	}

	// Use domain method to update
	err = blogPost.Update(title, content)
	if err != nil {
		return nil, err
	}

	err = u.Repo.Save(blogPost)
	if err != nil {
		u.Logger.Error("Failed to update blog post", "error", err, "id", blogPost.ID)
		return nil, err
	}
	return blogPost, nil
}

func (u *BlogPostUseCase) DeleteBlogPost(id string) error {
	err := u.Repo.Delete(id)
	if err != nil {
		u.Logger.Error("Failed to delete blog post", "error", err, "id", id)
		return err
	}
	return nil
}
