package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		switch {
		case !errors.Is(err, ErrNotFound):
			return nil, ErrNotFound
		}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:    req.Email,
		FullName: req.FullName,
		Phone:    req.Phone,
	}

	if err := s.repo.Create(ctx, user, string(hashed)); err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: "",
		User:  *user,
	}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	return nil, nil
}
