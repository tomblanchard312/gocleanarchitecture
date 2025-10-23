package sqlite

import (
	"database/sql"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteBlogPostRepository struct {
	DB *sql.DB
}

func NewSQLiteBlogPostRepository(db *sql.DB) interfaces.BlogPostRepository {
	return &SQLiteBlogPostRepository{DB: db}
}

func (r *SQLiteBlogPostRepository) Save(blogPost *entities.BlogPost) error {
	now := time.Now()
	if blogPost.CreatedAt.IsZero() {
		blogPost.CreatedAt = now
	}
	blogPost.UpdatedAt = now

	_, err := r.DB.Exec(`
        INSERT OR REPLACE INTO blog_posts (id, title, content, author_id, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `, blogPost.ID, blogPost.Title, blogPost.Content, blogPost.AuthorID, blogPost.CreatedAt, blogPost.UpdatedAt)
	return err
}

func (r *SQLiteBlogPostRepository) FindAll() ([]*entities.BlogPost, error) {
	rows, err := r.DB.Query("SELECT id, title, content, author_id, created_at, updated_at FROM blog_posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogPosts []*entities.BlogPost
	for rows.Next() {
		bp := &entities.BlogPost{}
		err := rows.Scan(&bp.ID, &bp.Title, &bp.Content, &bp.AuthorID, &bp.CreatedAt, &bp.UpdatedAt)
		if err != nil {
			return nil, err
		}
		blogPosts = append(blogPosts, bp)
	}
	return blogPosts, nil
}

func (r *SQLiteBlogPostRepository) FindByID(id string) (*entities.BlogPost, error) {
	bp := &entities.BlogPost{}
	err := r.DB.QueryRow("SELECT id, title, content, author_id, created_at, updated_at FROM blog_posts WHERE id = ?", id).
		Scan(&bp.ID, &bp.Title, &bp.Content, &bp.AuthorID, &bp.CreatedAt, &bp.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return bp, nil
}

func (r *SQLiteBlogPostRepository) Delete(id string) error {
	_, err := r.DB.Exec("DELETE FROM blog_posts WHERE id = ?", id)
	return err
}
