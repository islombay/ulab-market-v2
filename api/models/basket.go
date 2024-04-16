package models

import "time"

type BasketModel struct {
	UserID    string     `db:"user_id" json:"-"`
	ProductID string     `db:"product_id" json:"product_id"`
	Quantity  int        `db:"quantity" json:"quantity"`
	CreatedAt time.Time  `db:"created_at" json:"-"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}
