package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"victortillett.net/basic/internal/validator"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Comment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

type CommentModel struct {
	DB *sql.DB
}

func (c CommentModel) Insert(comment *Comment) error {
	query := `
		INSERT INTO comments (content, author)
		VALUES ($1, $2)
		RETURNING id, created_at, version`
	args := []any{comment.Content, comment.Author}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return c.DB.QueryRowContext(ctx, query, args...).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.Version,
	)
}

func (c CommentModel) Get(id int64) (*Comment, error) {
	query := `
		SELECT id, created_at, content, author, version
		FROM comments
		WHERE id = $1`
	var comment Comment
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.Content,
		&comment.Author,
		&comment.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &comment, nil
}

func (c CommentModel) Update(comment *Comment) error {
	query := `
		UPDATE comments
		SET content = $1, author = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version`
	args := []any{
		comment.Content,
		comment.Author,
		comment.ID,
		comment.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (c CommentModel) Delete(id int64) error {
	query := `DELETE FROM comments WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func ValidateComment(v *validator.Validator, comment *Comment) {
	v.Check(comment.Content != "", "content", "must be provided")
	v.Check(len(comment.Content) <= 100, "content", "must not be more than 100 bytes long")
	v.Check(comment.Author != "", "author", "must be provided")
	v.Check(len(comment.Author) <= 25, "author", "must not be more than 25 bytes long")
}
