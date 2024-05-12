package models_v1

import (
	"time"
)

// income
type Income struct {
	ID         string          `json:"id"`
	BranchID   string          `json:"branch_id"`
	TotalPrice *float32        `json:"total_price"`
	Comment    string          `json:"comment,omitempty"`
	Products   []IncomeProduct `json:"products"`
	CreatedAt  time.Time       `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt  *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`
}

type CreateIncome struct {
	BranchID string                `json:"branch_id" binding:"required"`
	Comment  string                `json:"comment"`
	Products []CreateIncomeProduct `json:"products" binding:"required"`
}

type IncomeRequest struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
}

type IncomeResponse struct {
	Incomes []Income `json:"incomes"`
	Count   int      `json:"count"`
}

// income_product
type IncomeProduct struct {
	ID           string     `json:"id"`
	IncomeID     string     `json:"income_id"`
	ProductID    string     `json:"product_id"`
	Quantity     int        `json:"quantity"`
	ProductPrice float32    `json:"product_price"`
	TotalPrice   float32    `json:"total_price"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}

type CreateIncomeProduct struct {
	IncomeID     string  `json:"income_id" swaggerignore:"true"`
	ProductID    string  `json:"product_id" binding:"required"`
	Quantity     int     `json:"quantity" binding:"required"`
	ProductPrice float32 `json:"product_price" binding:"required"`
}

type IncomeProductRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type IncomeProductResponse struct {
	IncomeProducts []IncomeProduct `json:"income_products"`
	Count          int             `json:"count"`
}

type CreateIncomeResponse struct {
	Income         Income          `json:"income"`
	IncomeProducts []IncomeProduct `json:"income_products"`
}
