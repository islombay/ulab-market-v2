package models

import "time"

type OrderModel struct {
	ID      string `db:"id" json:"id"`
	OrderID string `json:"order_id"`
	UserID  string `db:"user_id" json:"user_id"`

	ClientFirstName *string `json:"client_first_name"`
	ClientLastName  *string `json:"client_last_name"`
	ClientPhone     *string `json:"client_phone_number"`
	ClientComment   *string `json:"client_comment"`

	Status     string  `db:"status" json:"status"`
	TotalPrice float64 `db:"total_price" json:"total_price"`

	PaymentType string `db:"payment_type" json:"payment_type"`

	DeliveryType     string  `json:"delivery_type"`
	DeliveryAddrLat  float64 `json:"delivery_addr_lat"`
	DeliveryAddrLong float64 `json:"delivery_addr_long"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`

	Products []OrderProductModel `json:"products"`
}

type OrderProductModel struct {
	ID           string  `db:"id" json:"id"`
	OrderID      *string `db:"order_id" json:"order_id,omitempty"`
	ProductID    string  `db:"product_id" json:"product_id"`
	Quantity     int     `db:"quantity" json:"quantity"`
	ProductPrice float64 `db:"product_price" json:"product_price"`
	TotalPrice   float64 `db:"total_price" json:"total_price"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
