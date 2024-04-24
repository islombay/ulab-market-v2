package models_v1

type AddToFavourite struct {
	ProductID string `json:"product_id" binding:"required"`
}
