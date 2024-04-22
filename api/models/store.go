package models

import (
	"time"
)

type Storage struct {
	ID         string     `json:"id"`
	ProductID  string     `json:"product_id"`
	BranchID   string     `json:"branch_id"`
	TotalPrice float32    `json:"total_price"`
	Quantity   int        `json:"quantity"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" db:"deleted_at"`
}
