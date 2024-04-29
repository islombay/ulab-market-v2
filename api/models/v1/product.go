package models_v1

import (
	"app/api/models"
	"mime/multipart"
	"time"
)

type CreateProduct struct {
	Articul string `json:"articul" form:"articul" binding:"required"`
	NameUz  string `json:"name_uz" form:"name_uz" binding:"required"`
	NameRu  string `json:"name_ru" form:"name_ru" binding:"required"`

	DescriptionUz string `json:"description_uz" form:"description_uz" binding:"required"`
	DescriptionRu string `json:"description_ru" form:"description_ru" binding:"required"`

	OutcomePrice float64 `json:"outcome_price" form:"outcome_price" binding:"required"`

	IncomePrice float32 `json:"income_price" form:"income_price" binding:"required"`

	Quantity uint   `json:"quantity" form:"quantity" binding:"required"`
	BranchID string `json:"branch_id" form:"branch_id" binding:"required"`

	CategoryID string `json:"category_id" form:"category_id"`
	BrandID    string `json:"brand_id" form:"brand_id"`

	Status string `json:"status" form:"status"`

	MainImage *multipart.FileHeader `form:"main_image" binding:"required" swaggerignore:"true"`

	ImageFiles []*multipart.FileHeader `form:"image_files" binding:"required" swaggerignore:"true"`
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
	ID      string `json:"id" obj:"id"`
	Articul string `json:"articul,omitempty"`
	NameUz  string `json:"name_uz,omitempty" obj:"name_uz"`
	NameRu  string `json:"name_ru,omitempty" obj:"name_ru"`

	DescriptionUz string `json:"description_uz,omitempty" obj:"description_uz"`
	DescriptionRu string `json:"description_ru,omitempty" obj:"description_ru"`

	Price float64 `json:"price,omitempty" obj:"price"`

	Quantity int `json:"quantity,omitempty" obj:"quantity"`

	CategoryID string `json:"category_id,omitempty" obj:"category_id"`
	BrandID    string `json:"brand_id,omitempty" obj:"brand_id"`

	CategoryInformation models.Category `json:"category,omitempty" `
	// BrandInformation    models.Brand    `json:"brand_information,omitempty"`

	MainImage string  `json:"main_image,omitempty" obj:"main_image"`
	Rating    float32 `json:"rating,omitempty" obj:"rating"`

	ImageFiles []ProductMediaFiles `json:"image_files,omitempty" obj:"image_files"`
	VideoFiles []ProductMediaFiles `json:"video_files,omitempty" obj:"video_files"`

	CreatedAt time.Time `db:"created_at" json:"created_at" obj:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" obj:"updated_at"`
}

type ProductMediaFiles struct {
	ID        string `json:"id" obj:"id"`
	MediaFile string `json:"media_file" obj:"media_file"`
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

	// OutcomePrice float64 `json:"outcome_price"binding:"required"`

	CategoryID string `json:"category_id"`
	BrandID    string `json:"brand_id"`

	// Status string `json:"status" binding:"required"`
}

type ChangeProductPrice struct {
	ID    string  `json:"id" binding:"required"`
	Price float32 `json:"price" binding:"required"`
}
