package sqlite

import (
	"database/sql"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
)

type SQLiteCommentRepository struct {
	DB *sql.DB
}

func NewSQLiteCommentRepository(db *sql.DB) interfaces.CommentRepository {
	return &SQLiteCommentRepository{DB: db}
}

func (r *SQLiteCommentRepository) Save(comment *entities.Comment) error {
	_, err := r.DB.Exec(`
		INSERT OR REPLACE INTO comments (id, blog_post_id, author_id, content, parent_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, comment.ID, comment.BlogPostID, comment.AuthorID, comment.Content, comment.ParentID, comment.CreatedAt, comment.UpdatedAt)
	return err
}

func (r *SQLiteCommentRepository) FindByID(id string) (*entities.Comment, error) {
	comment := &entities.Comment{}
	err := r.DB.QueryRow(`
		SELECT id, blog_post_id, author_id, content, parent_id, created_at, updated_at
		FROM comments WHERE id = ?
	`, id).Scan(
		&comment.ID, &comment.BlogPostID, &comment.AuthorID, &comment.Content,
		&comment.ParentID, &comment.CreatedAt, &comment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return comment, nil
}

func (r *SQLiteCommentRepository) FindByBlogPostID(blogPostID string) ([]*entities.Comment, error) {
	rows, err := r.DB.Query(`
		SELECT id, blog_post_id, author_id, content, parent_id, created_at, updated_at
		FROM comments
		WHERE blog_post_id = ?
		ORDER BY created_at ASC
	`, blogPostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entities.Comment
	for rows.Next() {
		comment := &entities.Comment{}
		err := rows.Scan(
			&comment.ID, &comment.BlogPostID, &comment.AuthorID, &comment.Content,
			&comment.ParentID, &comment.CreatedAt, &comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

func (r *SQLiteCommentRepository) FindRepliesByParentID(parentID string) ([]*entities.Comment, error) {
	rows, err := r.DB.Query(`
		SELECT id, blog_post_id, author_id, content, parent_id, created_at, updated_at
		FROM comments
		WHERE parent_id = ?
		ORDER BY created_at ASC
	`, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entities.Comment
	for rows.Next() {
		comment := &entities.Comment{}
		err := rows.Scan(
			&comment.ID, &comment.BlogPostID, &comment.AuthorID, &comment.Content,
			&comment.ParentID, &comment.CreatedAt, &comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

func (r *SQLiteCommentRepository) Delete(id string) error {
	_, err := r.DB.Exec("DELETE FROM comments WHERE id = ?", id)
	return err
}

func (r *SQLiteCommentRepository) GetAll() ([]*entities.Comment, error) {
	rows, err := r.DB.Query(`
		SELECT id, blog_post_id, author_id, content, parent_id, created_at, updated_at
		FROM comments
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entities.Comment
	for rows.Next() {
		comment := &entities.Comment{}
		err := rows.Scan(
			&comment.ID, &comment.BlogPostID, &comment.AuthorID, &comment.Content,
			&comment.ParentID, &comment.CreatedAt, &comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}
