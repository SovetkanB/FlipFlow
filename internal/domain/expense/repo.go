package expense

import (
	"context"
	"fmt"

	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	List(ctx context.Context, projectID string) ([]Expense, error)
	Create(ctx context.Context, expense Expense) (*Expense, error)
	Delete(ctx context.Context, projectID, expenseID string) error
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) List(ctx context.Context, projectID string) ([]Expense, error) {
	query := `
		SELECT id, project_id, user_id, category, amount, description, paid_at, created_at, updated_at
		FROM expenses
		WHERE project_id = $1
		ORDER BY paid_at DESC, created_at DESC
	`

	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Expense
	for rows.Next() {
		var e Expense
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.UserID, &e.Category, &e.Amount, &e.Description,
			&e.PaidAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}

	return list, nil
}

func (r *repo) Create(ctx context.Context, expense Expense) (*Expense, error) {
	query := `
		INSERT INTO expenses (project_id, user_id, category, amount, description, paid_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx, query,
		expense.ProjectID, expense.UserID, expense.Category, expense.Amount, expense.Description, expense.PaidAt,
	).Scan(&expense.ID, &expense.CreatedAt, &expense.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("create expense: %w", err)
	}

	return &expense, nil
}

func (r *repo) Delete(ctx context.Context, userID, expenseID string) error {
	query := `
		DELETE FROM expenses
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.Exec(ctx, query, expenseID, userID)
	if err != nil {
		return fmt.Errorf("delete expense: %w", err)
	}

	if result.RowsAffected() == 0 {
		return response.ErrNotFound
	}

	return nil
}
