package project

import (
	"context"
	"errors"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, userID string, req CreateProjectRequest) (*Project, error) {
	p := Project{
		UserID:        userID,
		Title:         req.Title,
		Address:       req.Address,
		City:          req.City,
		AreaSqm:       req.AreaSqm,
		Rooms:         req.Rooms,
		Floor:         req.Floor,
		TotalFloors:   req.TotalFloors,
		PurchasePrice: req.PurchasePrice,
		TargetPrice:   req.TargetPrice,
		Description:   req.Description,
	}

	project, err := s.repo.Create(ctx, userID, p)
	if err != nil {
		return nil, err
	}

	return project, nil

}

func (s *Service) GetByID(ctx context.Context, projectID, userID string) (*Project, error) {
	project, err := s.repo.GetByID(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}

	total := s.repo.TotalExpenses(ctx, projectID)
	project.TotalExpenses = &total

	return project, nil
}

func (s Service) List(ctx context.Context, userID string, f ListFilter) ([]Project, int, error) {
	list, total, err := s.repo.List(ctx, userID, f)
	if err != nil {
		return nil, 0, err
	}

	for i := range list {
		t := s.repo.TotalExpenses(ctx, list[i].ID)
		list[i].TotalExpenses = &t
	}

	return list, total, nil
}

func (s *Service) Update(ctx context.Context, projectID, userID string, req UpdateProjectRequest) (*Project, error) {
	return s.repo.Update(ctx, projectID, userID, req)
}

func (s *Service) GetFinancialSummary(ctx context.Context, projectID, userID string) (*FinancialSummary, error) {
	project, err := s.repo.Summary(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}

	cats, _ := s.repo.GetExpensesByCategory(ctx, projectID)
	totalExpenses := s.repo.TotalExpenses(ctx, projectID)

	summary := &FinancialSummary{
		PurchasePrice:      project.PurchasePrice,
		TotalExpenses:      totalExpenses,
		TargetPrice:        project.TargetPrice,
		SoldPrice:          project.SoldPrice,
		ExpensesByCategory: cats,
	}

	invested := totalExpenses
	if project.PurchasePrice != nil {
		invested += *project.PurchasePrice
	}
	summary.TotalInvested = invested

	if project.TargetPrice != nil {
		profit := *project.TargetPrice - invested
		summary.EstimatedProfit = &profit
	}

	if project.SoldPrice != nil {
		profit := *project.SoldPrice - invested
		summary.ActualProfit = &profit
	}

	return summary, nil
}

func (s *Service) Delete(ctx context.Context, projectID, userID string) error {
	return s.repo.Delete(ctx, projectID, userID)
}

func (s *Service) ChangeStatus(ctx context.Context, projectID, userID string, newStatus Status, soldPrice *float64) (*Project, error) {
	current, err := s.repo.GetStatus(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}
	if !current.CanTransitionTo(newStatus) {
		return nil, errors.New("invalid status transition")
	}
	return s.repo.ChangeStatus(ctx, projectID, userID, newStatus, soldPrice)
}
