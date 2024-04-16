package models_v1

type AddToBasket struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
