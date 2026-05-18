package expense

import (
	"context"
	"time"
)

type Service interface {
	List(ctx context.Context, projectID string) ([]Expense, error)
	Create(ctx context.Context, projectID, userID string, req CreateExpenseRequest) (*Expense, error)
	Delete(ctx context.Context, projectID, expenseID string) error
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, projectID string) ([]Expense, error) {
	result, err := s.repo.List(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) Create(ctx context.Context, projectID, userID string, req CreateExpenseRequest) (*Expense, error) {
	if req.PaidAt == nil {
		now := time.Now()
		req.PaidAt = &now
	}

	expense := Expense{
		ProjectID:   projectID,
		UserID:      userID,
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
		PaidAt:      *req.PaidAt,
	}

	result, err := s.repo.Create(ctx, expense)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) Delete(ctx context.Context, projectID, expenseID string) error {
	return s.repo.Delete(ctx, projectID, expenseID)
}
