package models_v1

type CreateOrder struct {
	PaymentType string `json:"payment_type" binding:"required"`

	ClientFirstName string  `json:"client_first_name" binding:"required"`
	ClientLastName  string  `json:"client_last_name" binding:"required"`
	ClientPhone     string  `json:"client_phone_number" binding:"required"`
	ClientComment   *string `json:"client_comment"`

	DeliveryType     string  `json:"delivery_type" binding:"required"`
	DeliveryAddrLat  float64 `json:"delivery_addr_lat" binding:"required"`
	DeliveryAddrLong float64 `json:"delivery_addr_long" binding:"required"`
}
