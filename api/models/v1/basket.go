package models_v1

type AddToBasket struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  uint   `json:"quantity" binding:"required"`
}

type RemoveFromBasket struct {
	ProductID string `json:"product_id" binding:"required"`
}

type ChangeBasket struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  uint   `json:"quantity" binding:"required"`
}

type GetBasket struct {
	Products   []GetBasketProduct `json:"products"`
	TotalPrice float64            `json:"total_price"`
}

type GetBasketProduct struct {
	ID        string  `json:"id"`
	NameRu    string  `json:"name_ru"`
	NameUz    string  `json:"name_uz"`
	Price     float64 `json:"price"`
	MainImage *string `json:"main_image"`

	Quantity int `json:"quantity"`
}
