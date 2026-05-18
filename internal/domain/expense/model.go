package expense

import (
	"time"
)

type Category string

const (
	CategoryRoughWork       Category = "rough_work"
	CategoryFinishWork      Category = "finish_work"
	CategoryMaterialsRough  Category = "materials_rough"
	CategoryMaterialsFinish Category = "materials_finish"
	CategoryPlumbing        Category = "plumbing"
	CategoryElectrical      Category = "electrical"
	CategoryFurniture       Category = "furniture"
	CategoryElectronic      Category = "electronic"
	CategoryCommission      Category = "commission"
	CategoryTaxes           Category = "taxes"
	CategoryOther           Category = "other"
)

type Expense struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	UserID      string    `json:"user_id"`
	Category    Category  `json:"category"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description,omitempty"`
	PaidAt      time.Time `json:"paid_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateExpenseRequest struct {
	Category    Category   `json:"category" validate:"required"`
	Amount      float64    `json:"amount" validate:"required"`
	Description string     `json:"description" validate:"required"`
	PaidAt      *time.Time `json:"paid_at"`
}

type UpdateExpenseRequest struct {
	Category    *Category  `json:"category"`
	Amount      *float64   `json:"amount"`
	Description *string    `json:"description"`
	PaidAt      *time.Time `json:"paid_at"`
}
