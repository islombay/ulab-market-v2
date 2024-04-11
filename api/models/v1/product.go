package models_v1

import "mime/multipart"

type CreateProduct struct {
	Articul string `json:"articul" form:"articul" binding:"required"`
	NameUz  string `json:"name_uz" form:"name_uz"`
	NameRu  string `json:"name_ru" form:"name_ru" binding:"required"`

	DescriptionUz string `json:"description_uz" form:"description_uz" binding:"required"`
	DescriptionRu string `json:"description_ru" form:"description_ru" binding:"required"`

	IncomePrice  float32 `json:"income_price" form:"income_price"`
	OutcomePrice float64 `json:"outcome_price" form:"outcome_price" binding:"required"`

	Quantity int `json:"quantity" form:"quantity"`

	CategoryID string `json:"category_id" form:"category_id"`
	BrandID    string `json:"brand_id" form:"brand_id"`

	Status string `json:"status" form:"status" binding:"required"`

	MainImage *multipart.FileHeader `form:"main_image" swaggerignore:"true"`

	ImageFiles []*multipart.FileHeader `form:"image_files" swaggerignore:"true"`
	VideoFiles []*multipart.FileHeader `form:"video_files" swaggerignore:"true"`
}

type GetAllProductsQueryParams struct {
	CategoryID *string `form:"cid"`
	Query      *string `form:"q"`
	BrandID    *string `form:"bid"`
	Offset     int     `form:"offset"`
	Limit      int     `form:"limit"`
}

type ChangeProductMainImage struct {
	ProductID string                `form:"product_id" binding:"required"`
	Image     *multipart.FileHeader `form:"image" binding:"required" swaggerignore:"true"`
}
