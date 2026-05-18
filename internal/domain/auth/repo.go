package auth

import (
	"context"
	"strings"

	"github.com/SovetkanB/FlipFlow/internal/domain/user"
	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(ctx context.Context, user *user.User) error {
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

func (r *Repo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, phone, created_at, updated_at FROM users
		WHERE email = $1
	`
	var user user.User
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

func (r *Repo) GetByID(ctx context.Context, id string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, phone, created_at, updated_at FROM users
		WHERE id = $1
	`
	var user user.User
	err := r.db.QueryRow(ctx, query, id).Scan(
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

func (r *Repo) CreateRefreshToken(ctx context.Context, rt *RefreshToken) error {
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

func (r *Repo) DeleteRefreshToken(ctx context.Context, rt string) (string, error) {
	query := `
		DELETE FROM refresh_tokens
		WHERE token = $1 AND expires_at > NOW()
		RETURNING user_id
	`
	var userID string
	err := r.db.QueryRow(ctx, query, rt).Scan(&userID)

	if err != nil {
		return "", response.ErrNotValidRefreshToken
	}

	return userID, nil
}
