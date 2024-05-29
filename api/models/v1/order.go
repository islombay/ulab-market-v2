package models_v1

type CreateOrder struct {
	PaymentType     string `json:"payment_type" binding:"required"`
	PaymentCardType string `json:"payment_card_type" binding:"required"`

	ClientFirstName string  `json:"client_first_name" binding:"required"`
	ClientLastName  string  `json:"client_last_name" binding:"required"`
	ClientPhone     string  `json:"client_phone_number" binding:"required"`
	ClientComment   *string `json:"client_comment"`

	DeliveryType     string  `json:"delivery_type" binding:"required"`
	DeliveryAddrLat  float64 `json:"delivery_addr_lat" binding:"required"`
	DeliveryAddrLong float64 `json:"delivery_addr_long" binding:"required"`
	DeliverAddrName  *string `json:"delivery_name"`

	// Products []BasketProduct `json:"products" binding:"required"`
}

type BasketProduct struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  uint   `json:"quantity" binding:"required"`
}
