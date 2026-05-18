package expense

import (
	"context"
	"time"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, projectID, userID string, req CreateExpenseRequest) (*Expense, error) {
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

	return s.repo.Create(ctx, expense)
}

func (s *Service) List(ctx context.Context, projectID, userID string) ([]Expense, error) {
	if err := s.repo.OwnerCheck(ctx, projectID, userID); err != nil {
		return nil, err
	}

	return s.repo.List(ctx, projectID)
}

func (s *Service) Update(ctx context.Context, projectID, userID string, req UpdateExpenseRequest) (*Expense, error) {
	return s.repo.Update(ctx, projectID, userID, req)
}

func (s *Service) Delete(ctx context.Context, expenseID, userID string) error {
	return s.repo.Delete(ctx, userID, expenseID)
}
