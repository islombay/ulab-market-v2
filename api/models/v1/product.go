package models_v1

import (
	"app/api/models"
	"mime/multipart"
	"time"
)

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

type Product struct {
	ID     string `json:"id" obj:"id"`
	NameUz string `json:"name_uz" obj:"name_uz"`
	NameRu string `json:"name_ru" obj:"name_ru"`

	DescriptionUz string `json:"description_uz" obj:"description_uz"`
	DescriptionRu string `json:"description_ru" obj:"description_ru"`

	Price float64 `json:"price" obj:"price"`

	Quantity int `json:"quantity" obj:"quantity"`

	CategoryID string `json:"category_id" obj:"category_id"`
	BrandID    string `json:"brand_id" obj:"brand_id"`

	MainImage string  `json:"main_image" obj:"main_image"`
	Rating    float32 `json:"rating" obj:"rating"`

	ImageFiles []models.ProductMediaFiles `json:"image_files" obj:"image_files"`
	VideoFiles []models.ProductMediaFiles `json:"video_files" obj:"video_files"`

	CreatedAt time.Time `db:"created_at" json:"created_at" obj:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" obj:"updated_at"`
}

type AddProductMediaFiles struct {
	ProductID  string                  `form:"product_id" binding:"required"`
	MediaFiles []*multipart.FileHeader `form:"media_files" swaggerignore:"true" binding:"required"`
}

type ChangeProductRequest struct {
	ID      string `json:"id" binding:"required"`
	Articul string `json:"articul" binding:"required"`
	NameUz  string `json:"name_uz"`
	NameRu  string `json:"name_ru"`

	DescriptionUz string `json:"description_uz" binding:"required"`
	DescriptionRu string `json:"description_ru" binding:"required"`

	IncomePrice  float32 `json:"income_price"`
	OutcomePrice float64 `json:"outcome_price"binding:"required"`

	Quantity int `json:"quantity"`

	CategoryID string `json:"category_id"`
	BrandID    string `json:"brand_id"`

	Status string `json:"status" binding:"required"`
}
