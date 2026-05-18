package project

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	Create(ctx context.Context, userID string, project Project) (*Project, error)
	Update(ctx context.Context, sets []string, args []interface{}, argN int) (*Project, error)
	Delete(ctx context.Context, projectID, userID string) error
	GetByID(ctx context.Context, projectID, userID string) (*Project, error)
	List(ctx context.Context, args []interface{}, whereClause string, limit, offset int) ([]Project, int, error)
	GetExpensesByCategory(ctx context.Context, propertyID string) ([]CategorySum, error)
	GetSummary(ctx context.Context, projectID, userID string) (*Project, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{db: db}
}

func (r *repo) GetSummary(ctx context.Context, projectID, userID string) (*Project, error) {
	query := `
		SELECT purchase_price, target_price, sold_price
		FROM projects WHERE id = $1 AND user_id = $2
	`
	var project Project
	err := r.db.QueryRow(ctx, query, projectID, userID).Scan(&project.PurchasePrice, &project.TargetPrice, &project.SoldPrice)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.ErrNotFound
		}
		return nil, fmt.Errorf("get property: %w", err)
	}

	return &project, nil
}

func (r *repo) Create(ctx context.Context, userID string, project Project) (*Project, error) {
	query := `
		INSERT INTO projects (user_id, title, address, city, area_sqm, rooms,
		purchase_price, target_price, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, status, created_at, updated_at
	`
	err := r.db.QueryRow(
		ctx, query,
		userID, project.Title, project.Address, project.City, project.AreaSqm, project.Rooms,
		project.PurchasePrice, project.TargetPrice, project.Description,
	).Scan(
		&project.ID, &project.Status, &project.CreatedAt, &project.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create property: %w", err)
	}

	return &project, nil
}

func (r *repo) Delete(ctx context.Context, projectID, userID string) error {
	query := `
		DELETE FROM projects
		WHERE id = $1 AND user_id = $2
	`
	result, err := r.db.Exec(ctx, query, projectID, userID)

	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}

	if result.RowsAffected() == 0 {
		return response.ErrNotFound
	}
	return nil
}

func (r *repo) GetByID(ctx context.Context, projectID, userID string) (*Project, error) {
	query := `
		SELECT id, user_id, title, address, city, area_sqm, rooms,
				purchase_price, target_price, sold_price, status, description,
				purchased_at, sold_at, created_at, updated_at
		FROM projects
		WHERE id = $1 AND user_id = $2
	`

	var p Project
	err := r.db.QueryRow(ctx, query, projectID, userID).Scan(
		&p.ID, &p.UserID, &p.Title, &p.Address, &p.City, &p.AreaSqm, &p.Rooms,
		&p.PurchasePrice, &p.TargetPrice, &p.SoldPrice, &p.Status, &p.Description,
		&p.PurchasedAt, &p.SoldAt, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.ErrNotFound
		}
		return nil, fmt.Errorf("get property: %w", err)
	}

	return &p, nil
}

func (r *repo) Update(ctx context.Context, sets []string, args []interface{}, argN int) (*Project, error) {
	query := fmt.Sprintf(`
		UPDATE projects
		SET %s , updated_at = NOW()
		WHERE id = $%d AND user_id = $%d
		RETURNING id, user_id, title, address, city, area_sqm, rooms,
		        	purchase_price, target_price,
		          sold_price, status, description, purchased_at, sold_at,
		          created_at, updated_at
	`, strings.Join(sets, ", "), argN, argN+1)

	var p Project
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.UserID, &p.Title, &p.Address, &p.City, &p.AreaSqm,
		&p.Rooms, &p.PurchasePrice, &p.TargetPrice,
		&p.SoldPrice, &p.Status, &p.Description, &p.PurchasedAt, &p.SoldAt,
		&p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, response.ErrNotFound
		default:
			return nil, fmt.Errorf("update property: %w", err)
		}
	}
	return &p, nil
}

func (r *repo) List(ctx context.Context, args []interface{}, whereClause string, limit, offset int) ([]Project, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		"SELECT COUNT(*) FROM projects "+whereClause, args...,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, fmt.Sprintf(`
		SELECT id, user_id, title, address, city, area_sqm, rooms,
		       purchase_price, target_price,
		       sold_price, status, description, purchased_at, sold_at,
		       created_at, updated_at
		FROM projects
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.Title, &p.Address, &p.City, &p.AreaSqm,
			&p.Rooms, &p.PurchasePrice, &p.TargetPrice,
			&p.SoldPrice, &p.Status, &p.Description, &p.PurchasedAt, &p.SoldAt,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, p)
	}

	return list, total, nil
}

func (r *repo) GetExpensesByCategory(ctx context.Context, propertyID string) ([]CategorySum, error) {
	query := `
		SELECT category, SUM(amount)
		FROM expenses
		WHERE project_id = $1
		GROUP BY category
		ORDER BY category
	`

	rows, err := r.db.Query(ctx, query, propertyID)
	if err != nil {
		return nil, err
	}

	var list []CategorySum
	for rows.Next() {
		var c CategorySum
		if err := rows.Scan(&c.Category, &c.Total); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}
