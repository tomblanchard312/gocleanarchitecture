package usecases

import (
	"gocleanarchitecture/entities"
	"gocleanarchitecture/frameworks/logger"
	"gocleanarchitecture/interfaces"
)

type BlogPostUseCase struct {
	Repo   interfaces.BlogPostRepository
	Logger logger.Logger
}

func (u *BlogPostUseCase) CreateBlogPost(blogPost *entities.BlogPost) error {
	err := u.Repo.Save(blogPost)
	if err != nil {
		u.Logger.Error("Failed to create blog post", logger.Field("error", err))
		return err
	}
	return nil
}

func (u *BlogPostUseCase) GetAllBlogPosts() ([]*entities.BlogPost, error) {
	blogPosts, err := u.Repo.FindAll()
	if err != nil {
		u.Logger.Error("Failed to get all blog posts", logger.Field("error", err))
		return nil, err
	}
	return blogPosts, nil
}

func (u *BlogPostUseCase) GetBlogPost(id string) (*entities.BlogPost, error) {
	blogPost, err := u.Repo.FindByID(id)
	if err != nil {
		u.Logger.Error("Failed to get blog post", logger.Field("error", err), logger.Field("id", id))
		return nil, err
	}
	return blogPost, nil
}

func (u *BlogPostUseCase) UpdateBlogPost(blogPost *entities.BlogPost) error {
	err := u.Repo.Save(blogPost)
	if err != nil {
		u.Logger.Error("Failed to update blog post", logger.Field("error", err), logger.Field("id", blogPost.ID))
		return err
	}
	return nil
}

func (u *BlogPostUseCase) DeleteBlogPost(id string) error {
	err := u.Repo.Delete(id)
	if err != nil {
		u.Logger.Error("Failed to delete blog post", logger.Field("error", err), logger.Field("id", id))
		return err
	}
	return nil
}
