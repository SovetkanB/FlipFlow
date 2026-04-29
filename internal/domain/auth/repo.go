package auth

import (
	"context"
	"strings"

	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	CreateRefreshToken(ctx context.Context, rt *RefreshToken) error
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Phone,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "unique"):
			return response.ErrEmailTaken
		default:
			return err
		}
	}

	return nil
}

func (r *repo) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, full_name, phone, created_at, updated_at FROM users
		WHERE email = $1
	`
	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "no rows"):
			return nil, response.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r *repo) CreateRefreshToken(ctx context.Context, rt *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query, rt.UserID, rt.Token, rt.ExpiresAt).Scan(
		&rt.ID,
		&rt.CreatedAt,
	)

	return err
}
