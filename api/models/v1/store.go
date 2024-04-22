package models_v1

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

type CreateStorage struct {
	ProductID  string  `json:"product_id"`
	BranchID   string  `json:"branch_id"`
	TotalPrice float32 `json:"total_price"`
	Quantity   int     `json:"quantity"`
}

type UpdateStorage struct {
	ID         string  `json:"id"`
	ProductID  string  `json:"product_id"`
	BranchID   string  `json:"branch_id"`
	TotalPrice float32 `json:"total_price"`
	Quantity   int     `json:"quantity"`
}

type StorageRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type StorageResponse struct {
	Storage []Storage `json:"storage"`
	Count   int       `json:"count"`
}
