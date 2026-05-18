package project

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Service interface {
	Create(ctx context.Context, userID string, req CreateProjectRequest) (*Project, error)
	Delete(ctx context.Context, projectID, userID string) error
	GetByID(ctx context.Context, projectID, userID string) (*Project, error)
	Update(ctx context.Context, projectID, userID string, req UpdateProjectRequest) (*Project, error)
	List(ctx context.Context, userID string, f ListFilter) ([]Project, int, error)
	GetFinancialSummary(ctx context.Context, projectID, userID string) (*FinancialSummary, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetFinancialSummary(ctx context.Context, projectID, userID string) (*FinancialSummary, error) {
	project, err := s.repo.GetSummary(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}

	cats, _ := s.repo.GetExpensesByCategory(ctx, projectID)
	totalExpenses := *sumCategories(cats)

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

func (s *service) Create(ctx context.Context, userID string, req CreateProjectRequest) (*Project, error) {
	p := Project{
		UserID:        userID,
		Title:         req.Title,
		Address:       req.Address,
		City:          req.City,
		AreaSqm:       req.AreaSqm,
		Rooms:         req.Rooms,
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

func (s *service) Delete(ctx context.Context, projectID, userID string) error {
	return s.repo.Delete(ctx, projectID, userID)
}

func (s *service) GetByID(ctx context.Context, projectID, userID string) (*Project, error) {
	project, err := s.repo.GetByID(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}

	project.ExpensesByCategory, _ = s.repo.GetExpensesByCategory(ctx, projectID)
	project.TotalExpenses = sumCategories(project.ExpensesByCategory)

	return project, nil
}

func (s *service) Update(ctx context.Context, projectID, userID string, req UpdateProjectRequest) (*Project, error) {
	sets := []string{}
	args := []interface{}{}
	argN := 1

	add := func(col string, value interface{}) {
		sets = append(sets, fmt.Sprintf("%s = $%d", col, argN))
		args = append(args, value)
		argN++
	}

	if req.Title != nil {
		add("title", *req.Title)
	}
	if req.Address != nil {
		add("address", *req.Address)
	}
	if req.City != nil {
		add("city", *req.City)
	}
	if req.AreaSqm != nil {
		add("area_sqm", *req.AreaSqm)
	}
	if req.Rooms != nil {
		add("rooms", *req.Rooms)
	}
	if req.PurchasePrice != nil {
		add("purchase_price", *req.PurchasePrice)
	}
	if req.TargetPrice != nil {
		add("target_price", *req.TargetPrice)
	}
	if req.Description != nil {
		add("description", *req.Description)
	}
	if req.Status != nil {
		add("status", *req.Status)

		if *req.Status == StatusPurchased {
			add("purchased_at", time.Now())
		}
		if *req.Status == StatusSold {
			if req.SoldPrice != nil {
				add("sold_price", *req.SoldPrice)
			}
			add("sold_at", time.Now())
		}
	}

	if len(sets) == 0 {
		return s.GetByID(ctx, projectID, userID)
	}

	args = append(args, projectID, userID)

	project, err := s.repo.Update(ctx, sets, args, argN)
	if err != nil {
		return nil, err
	}

	// TODO Пересчёт средних при продаже объекта

	project.ExpensesByCategory, _ = s.repo.GetExpensesByCategory(ctx, projectID)
	project.TotalExpenses = sumCategories(project.ExpensesByCategory)
	return project, nil
}

func (s service) List(ctx context.Context, userID string, f ListFilter) ([]Project, int, error) {
	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Offset < 0 {
		f.Limit = 0
	}

	args := []interface{}{userID}
	where := []string{"user_id = $1"}

	if f.Status != "" {
		args = append(args, f.Status)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	if f.City != "" {
		args = append(args, "%"+f.City+"%")
		where = append(where, fmt.Sprintf("city ILIKE $%d", len(args)))
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	list, total, err := s.repo.List(ctx, args, whereClause, f.Limit, f.Offset)
	if err != nil {
		return nil, 0, err
	}

	for i := range list {
		cats, _ := s.repo.GetExpensesByCategory(ctx, list[i].ID)
		list[i].TotalExpenses = sumCategories(cats)
	}

	return list, total, nil
}

func sumCategories(cats []CategorySum) *float64 {
	if len(cats) == 0 {
		return nil
	}
	var total float64
	for _, c := range cats {
		total += c.Total
	}
	return &total
}
