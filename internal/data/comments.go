package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Comment struct maps to JSON and DB
type Comment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

var ErrRecordNotFound = errors.New("record not found")

// CommentModel wraps DB access
type CommentModel struct {
	DB *sql.DB
}

// Insert new comment
func (m CommentModel) Insert(comment *Comment) error {
	query := `
		INSERT INTO comments (content, author)
		VALUES ($1, $2)
		RETURNING id, created_at, version`
	args := []any{comment.Content, comment.Author}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.Version,
	)
}

// Get comment by ID
func (m CommentModel) Get(id int64) (*Comment, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, content, author, version FROM comments WHERE id=$1`

	var c Comment
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.CreatedAt, &c.Content, &c.Author, &c.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &c, nil
}
