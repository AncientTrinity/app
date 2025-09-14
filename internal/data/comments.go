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

func ValidateComment(v *validator.Validator, comment *Comment) {
	v.Check(comment.Content != "", "content", "must be provided")
	v.Check(len(comment.Content) <= 100, "content", "must not be more than 100 bytes long")
	v.Check(comment.Author != "", "author", "must be provided")
	v.Check(len(comment.Author) <= 25, "author", "must not be more than 25 bytes long")
}
