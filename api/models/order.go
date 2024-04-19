package models

import "time"

type OrderModel struct {
	ID          string  `db:"id" json:"id"`
	UserID      string  `db:"user_id" json:"user_id"`
	Status      string  `db:"status" json:"status"`
	TotalPrice  float64 `db:"total_price" json:"total_price"`
	PaymentType string  `db:"payment_type" json:"payment_type"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type OrderProductModel struct {
	ID           string  `db:"id"`
	OrderID      string  `db:"order_id"`
	ProductID    string  `db:"product_id"`
	Quantity     int     `db:"quantity"`
	ProductPrice float64 `db:"product_price"`
	TotalPrice   float64 `db:"total_price"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
