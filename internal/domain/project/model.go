package project

import "time"

type Status string

const (
	StatusSearch     Status = "search"
	StatusPurchased  Status = "purchased"
	StatusRenovation Status = "renovation"
	StatusForSale    Status = "for_sale"
	StatusSold       Status = "sold"
)

type Project struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	Title         string     `json:"title"`
	Address       string     `json:"address"`
	City          string     `json:"city"`
	AreaSqm       *float64   `json:"area_sqm"`
	Rooms         *int       `json:"rooms"`
	PurchasePrice *float64   `json:"purchase_price"`
	TargetPrice   *float64   `json:"target_price"`
	SoldPrice     *float64   `json:"sold_price"`
	Status        Status     `json:"status"`
	Description   string     `json:"description"`
	PurchasedAt   *time.Time `json:"purchased_at"`
	SoldAt        *time.Time `json:"sold_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	TotalExpenses      *float64      `json:"total_expenses"`
	ExpensesByCategory []CategorySum `json:"expenses_by_category"`
}

type CategorySum struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

type FinancialSummary struct {
	PurchasePrice      *float64      `json:"purchase_price"`
	TotalExpenses      float64       `json:"total_expenses"`
	TotalInvested      float64       `json:"total_invested"`
	TargetPrice        *float64      `json:"target_price"`
	SoldPrice          *float64      `json:"sold_price"`
	EstimatedProfit    *float64      `json:"estimated_profit"`
	ActualProfit       *float64      `json:"actual_profit"`
	ExpensesByCategory []CategorySum `json:"expenses_by_category"`
}

type CreateProjectRequest struct {
	Title         string   `json:"title" validate:"required,min=2,max=255"`
	Address       string   `json:"address"`
	City          string   `json:"city"`
	AreaSqm       *float64 `json:"area_sqm"`
	Rooms         *int     `json:"rooms"`
	PurchasePrice *float64 `json:"purchase_price"`
	TargetPrice   *float64 `json:"target_price"`
	Description   string   `json:"description"`
}

type UpdateProjectRequest struct {
	Title         *string  `json:"title"`
	Address       *string  `json:"address"`
	City          *string  `json:"city"`
	AreaSqm       *float64 `json:"area_sqm"`
	Rooms         *int     `json:"rooms"`
	PurchasePrice *float64 `json:"purchase_price"`
	TargetPrice   *float64 `json:"target_price"`
	SoldPrice     *float64 `json:"sold_price"`
	Status        *Status  `json:"status"`
	Description   *string  `json:"description"`
}

type ListFilter struct {
	Status string
	City   string
	Limit  int
	Offset int
}
