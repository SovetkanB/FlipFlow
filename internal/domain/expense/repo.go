package expense

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) OwnerCheck(ctx context.Context, projectID, userID string) error {
	query := `
		SELECT user_id FROM projects WHERE id = $1
	`
	var ownerID string
	err := r.db.QueryRow(ctx, query, projectID).Scan(&ownerID)
	if err != nil || ownerID != userID {
		return response.ErrNotFound
	}
	return nil
}

func (r *Repo) Create(ctx context.Context, expense Expense) (*Expense, error) {
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

func (r *Repo) List(ctx context.Context, projectID string) ([]Expense, error) {
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

func (r *Repo) Update(ctx context.Context, projectID, userID string, expense UpdateExpenseRequest) (*Expense, error) {
	sets := []string{}
	args := []any{}
	argN := 1
	add := func(col string, value interface{}) {
		sets = append(sets, fmt.Sprintf("%s = %d", col, argN))
		args = append(args, value)
		argN++
	}
	if expense.Category != nil {
		add("category", *expense.Category)
	}
	if expense.Amount != nil {
		add("amount", *expense.Amount)
	}
	if expense.Description != nil {
		add("description", *expense.Description)
	}
	if expense.PaidAt != nil {
		add("paid_at", *expense.PaidAt)
	}

	if len(sets) == 0 {
		var e Expense
		err := r.db.QueryRow(ctx,
			`SELECT id, project_id, user_id, category, amount, 
					description, paid_at, created_at, updated_at
			 FROM expenses
			 WHERE id=$1 AND user_id=$2
			`, projectID, userID,
		).Scan(&e.ID, &e.ProjectID, &e.UserID, &e.Category, &e.Amount, &e.Description,
			&e.PaidAt, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				return nil, response.ErrNotFound
			default:
				return nil, fmt.Errorf("update expense: %w", err)
			}
		}
		return &e, nil
	}

	args = append(args, projectID, userID)
	query := fmt.Sprintf(`
		UPDATE expenses
		SET %s, updated_at = NOW
		WHERE id = $%d AND user_id $%d
		RETURNING id, project_id, user_id, category, amount, 
					description, paid_at, created_at, updated_at
	`, strings.Join(sets, ", "), argN, argN+1)

	var e Expense
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&e.ID, &e.ProjectID, &e.UserID, &e.Category, &e.Amount, &e.Description,
		&e.PaidAt, &e.CreatedAt, &e.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, response.ErrNotFound
		default:
			return nil, fmt.Errorf("update project: %w", err)
		}
	}

	return &e, nil
}

func (r *Repo) Delete(ctx context.Context, userID, expenseID string) error {
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
