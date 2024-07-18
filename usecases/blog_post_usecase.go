package usecases

import (
	"gocleanarchitecture/entities"
	"gocleanarchitecture/errors"
	"gocleanarchitecture/frameworks/logger"
)

type BlogPostRepository interface {
	Save(blogPost entities.BlogPost) error
	FindAll() ([]entities.BlogPost, error)
}

type BlogPostUseCaseInterface interface {
	CreateBlogPost(blogPost entities.BlogPost) error
	GetAllBlogPosts() ([]entities.BlogPost, error)
}

type BlogPostUseCase struct {
	Repo   BlogPostRepository
	Logger logger.Logger
}

func (uc *BlogPostUseCase) CreateBlogPost(blogPost entities.BlogPost) error {
	err := uc.Repo.Save(blogPost)
	if err != nil {
		uc.Logger.Error("Failed to create blog post", logger.Field("error", err))
		return errors.Wrap(err, "failed to create blog post")
	}
	uc.Logger.Info("Blog post created successfully", logger.Field("id", blogPost.ID))
	return nil
}

func (uc *BlogPostUseCase) GetAllBlogPosts() ([]entities.BlogPost, error) {
	posts, err := uc.Repo.FindAll()
	if err != nil {
		uc.Logger.Error("Failed to get all blog posts", logger.Field("error", err))
		return nil, errors.Wrap(err, "failed to get all blog posts")
	}
	uc.Logger.Info("Retrieved all blog posts", logger.Field("count", len(posts)))
	return posts, nil
}
