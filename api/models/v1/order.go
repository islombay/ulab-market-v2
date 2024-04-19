package models_v1

type CreateOrder struct {
	PaymentType string `json:"payment_type" binding:"required"`
}
