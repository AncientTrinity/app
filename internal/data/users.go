package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"-"` // don’t expose in JSON
	Activated    bool      `json:"activated"`
	Version      int       `json:"version"`
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash, activated, version)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Activated,
		user.Version,
	).Scan(&user.ID, &user.CreatedAt)
}

func (m UserModel) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, errors.New("invalid id")
	}

	query := `
		SELECT id, created_at, username, email, password_hash, activated, version
		FROM users
		WHERE id = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	return &user, nil
}
