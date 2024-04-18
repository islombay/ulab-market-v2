package models_v1

type AddToBasket struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type RemoveFromBasket struct {
	ProductID string `json:"product_id" binding:"required"`
}

type ChangeBasket struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}
