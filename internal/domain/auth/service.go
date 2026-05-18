package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/SovetkanB/FlipFlow/internal/config"
	"github.com/SovetkanB/FlipFlow/internal/domain/user"
	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repo
	cfg  config.JWTConfig
}

func NewService(repo *Repo, cfg config.JWTConfig) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &user.User{
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

func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil || !(bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) == nil) {
		return nil, response.ErrInvalidPassword
	}

	return s.issueTokens(ctx, user)
}

func (s *Service) Refresh(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error) {
	hash := sha256.Sum256([]byte(req.RefreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	userID, err := s.repo.DeleteRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user)
}

func (s *Service) Me(ctx context.Context, userID string) (*user.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *Service) ValidateToken(token string) (*Claims, error) {
	return parseToken(s.cfg.Secret, token)
}

func (s *Service) issueTokens(ctx context.Context, user *user.User) (*AuthResponse, error) {
	accessToken, err := generateAccessToken(s.cfg.Secret, s.cfg.AccessTTL, user)
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
