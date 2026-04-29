package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/SovetkanB/FlipFlow/internal/config"
	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}

type service struct {
	repo Repo
	cfg  config.JWTConfig
}

func NewService(repo Repo, cfg config.JWTConfig) Service {
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        req.Email,
		FullName:     req.FullName,
		PasswordHash: string(hashed),
		Phone:        req.Phone,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user)
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil || !(bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) == nil) {
		return nil, response.ErrInvalidPassword
	}

	return s.issueTokens(ctx, user)
}

func (s *service) issueTokens(ctx context.Context, user *User) (*AuthResponse, error) {
	accessToken, err := GenerateAccessToken(s.cfg, user)
	if err != nil {
		return nil, err
	}

	rawToken := make([]byte, 32)
	if _, err := rand.Read(rawToken); err != nil {
		return nil, err
	}
	refreshToken := hex.EncodeToString(rawToken)

	rt := &RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.cfg.RefreshTTL),
	}

	if err := s.repo.CreateRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: TokenPair{
			AccessToken:  accessToken,
			RefreshToken: rt.Token,
			ExpiresAt:    rt.ExpiresAt,
		},
		User: user.ToResponse(),
	}, nil
}
