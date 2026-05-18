package project

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, userID string, project Project) (*Project, error) {
	query := `
		INSERT INTO projects (user_id, title, address, city, area_sqm, rooms, floor, total_floor,
		purchase_price, target_price, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, status, created_at, updated_at
	`
	err := r.db.QueryRow(
		ctx, query,
		userID, project.Title, project.Address, project.City, project.AreaSqm, project.Rooms, project.Floor, project.TotalFloors,
		project.PurchasePrice, project.TargetPrice, project.Description,
	).Scan(
		&project.ID, &project.Status, &project.CreatedAt, &project.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	return &project, nil
}

func (r *Repo) GetByID(ctx context.Context, projectID, userID string) (*Project, error) {
	query := `
		SELECT id, user_id, title, address, city, area_sqm, rooms, floor, total_floor,
				purchase_price, target_price, sold_price, status, description,
				purchased_at, sold_at, created_at, updated_at
		FROM projects
		WHERE id = $1 AND user_id = $2
	`

	var p Project
	err := r.db.QueryRow(ctx, query, projectID, userID).Scan(
		&p.ID, &p.UserID, &p.Title, &p.Address, &p.City, &p.AreaSqm, &p.Rooms, &p.Floor, &p.TotalFloors,
		&p.PurchasePrice, &p.TargetPrice, &p.SoldPrice, &p.Status, &p.Description,
		&p.PurchasedAt, &p.SoldAt, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.ErrNotFound
		}
		return nil, fmt.Errorf("get project: %w", err)
	}

	return &p, nil
}

func (r *Repo) List(ctx context.Context, userID string, f ListFilter) ([]Project, int, error) {
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

	whereClause := "WHERE " + strings.Join(where, " AND ")

	var total int
	err := r.db.QueryRow(ctx,
		"SELECT COUNT(*) FROM projects "+whereClause, args...,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, f.Limit, f.Offset)
	rows, err := r.db.Query(ctx, fmt.Sprintf(`
		SELECT id, user_id, title, address, city, area_sqm, rooms, floor, total_floors,
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
			&p.Rooms, &p.Floor, &p.TotalFloors, &p.PurchasePrice, &p.TargetPrice,
			&p.SoldPrice, &p.Status, &p.Description, &p.PurchasedAt, &p.SoldAt,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, p)
	}

	return list, total, nil
}

func (r *Repo) Update(ctx context.Context, projectID, userID string, project UpdateProjectRequest) (*Project, error) {
	sets := []string{}
	args := []any{}
	argN := 1

	add := func(col string, value interface{}) {
		sets = append(sets, fmt.Sprintf("%s = $%d", col, argN))
		args = append(args, value)
		argN++
	}

	if project.Title != nil {
		add("title", *project.Title)
	}
	if project.Address != nil {
		add("address", *project.Address)
	}
	if project.City != nil {
		add("city", *project.City)
	}
	if project.AreaSqm != nil {
		add("area_sqm", *project.AreaSqm)
	}
	if project.Rooms != nil {
		add("rooms", *project.Rooms)
	}
	if project.Floor != nil {
		add("floor", *project.Floor)
	}
	if project.TotalFloors != nil {
		add("total_floors", *project.TotalFloors)
	}
	if project.PurchasePrice != nil {
		add("purchase_price", *project.PurchasePrice)
	}
	if project.TargetPrice != nil {
		add("target_price", *project.TargetPrice)
	}
	if project.Description != nil {
		add("description", *project.Description)
	}

	if len(sets) == 0 {
		return r.GetByID(ctx, projectID, userID)
	}
	args = append(args, projectID, userID)

	query := fmt.Sprintf(`
		UPDATE projects
		SET %s , updated_at = NOW()
		WHERE id = $%d AND user_id = $%d
		RETURNING id, user_id, title, address, city, area_sqm, rooms, floors, total_floors,
		        	purchase_price, target_price,
		          sold_price, status, description, purchased_at, sold_at,
		          created_at, updated_at
	`, strings.Join(sets, ", "), argN, argN+1)

	var p Project
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.UserID, &p.Title, &p.Address, &p.City, &p.AreaSqm,
		&p.Rooms, &p.Floor, &p.TotalFloors, &p.PurchasePrice, &p.TargetPrice,
		&p.SoldPrice, &p.Status, &p.Description, &p.PurchasedAt, &p.SoldAt,
		&p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, response.ErrNotFound
		default:
			return nil, fmt.Errorf("update project: %w", err)
		}
	}
	return &p, nil
}

func (r *Repo) ChangeStatus(ctx context.Context, projectID, userID string, newStatus Status, soldPrice *float64) (*Project, error) {
	sets := "status=$1"
	args := []any{newStatus}
	n := 2
	if newStatus == StatusPurchased {
		sets += fmt.Sprintf(", purchased_at=$%d", n)
		args = append(args, time.Now())
		n++
	}
	if newStatus == StatusSold {
		sets += fmt.Sprintf(", sold_at=$%d", n)
		args = append(args, time.Now())
		n++
		if soldPrice != nil {
			sets += fmt.Sprintf(", sold_price=$%d", n)
			args = append(args, *soldPrice)
			n++
		}
	}
	args = append(args, projectID, userID)
	query := fmt.Sprintf(`
		UPDATE projects SET %s WHERE id=$%d AND user_id=$%d RETURNING id, user_id, title, address, city, area_sqm, rooms, floors, total_floors,
		        	purchase_price, target_price,
		          sold_price, status, description, purchased_at, sold_at,
		          created_at, updated_at
	`, sets, n, n+1)

	var p Project
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.UserID, &p.Title, &p.Address, &p.City, &p.AreaSqm,
		&p.Rooms, &p.Floor, &p.TotalFloors, &p.PurchasePrice, &p.TargetPrice,
		&p.SoldPrice, &p.Status, &p.Description, &p.PurchasedAt, &p.SoldAt,
		&p.CreatedAt, &p.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, response.ErrNotFound
	}
	return &p, err
}

func (r *Repo) GetStatus(ctx context.Context, projectID, userID string) (Status, error) {
	var s Status
	err := r.db.QueryRow(ctx, "SELECT status FROM projects WHERE id=$1 AND user_id=$2", projectID, userID).Scan(&s)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", response.ErrNotFound
	}
	return s, err
}

func (r *Repo) Delete(ctx context.Context, projectID, userID string) error {
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

func (r *Repo) Summary(ctx context.Context, projectID, userID string) (*Project, error) {
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
		return nil, fmt.Errorf("get project summary: %w", err)
	}

	return &project, nil
}

func (r *Repo) GetExpensesByCategory(ctx context.Context, propertyID string) ([]CategorySum, error) {
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

func (r *Repo) TotalExpenses(ctx context.Context, id string) float64 {
	var total float64
	r.db.QueryRow(ctx, "SELECT COALESCE(SUM(amount),0) FROM expenses WHERE project_id=$1", id).Scan(&total)
	return total
}
