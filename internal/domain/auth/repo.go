package auth

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	Create(ctx context.Context, user *User, passHash string) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, user *User, passHash string) error {
	query := `
		INSERT INTO users (email, password_hash, full_name, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Email,
		passHash,
		user.FullName,
		user.Phone,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return err
}

func (r *repo) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, full_name, phone, created_at, updated_at FROM users
		WHERE email = $1
	`
	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "no rows"):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
