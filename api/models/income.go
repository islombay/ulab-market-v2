package models

import "time"

type Income struct {
	ID         string     `json:"id"`
	BranchID   string     `json:"branch_id"`
	TotalPrice float32    `json:"total_price"`
	Comment    string     `json:"comment"`
	CourierID  string     `json:"courier_id"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" db:"deleted_at"`
}

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
