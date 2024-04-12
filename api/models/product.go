package models

import (
	"time"
)

type Product struct {
	ID      string `db:"id" json:"id" obj:"id"`
	Articul string `db:"articul" json:"articul,omitempty"`
	NameUz  string `db:"name_uz" json:"name_uz,omitempty" obj:"name_uz"`
	NameRu  string `db:"name_ru" json:"name_ru,omitempty" obj:"name_ru"`

	DescriptionUz string `db:"description_uz" json:"description_uz,omitempty" obj:"description_uz"`
	DescriptionRu string `db:"description_ru" json:"description_ru,omitempty" obj:"description_ru"`

	IncomePrice  float32 `db:"income_price" json:"income_price,omitempty"`
	OutcomePrice float64 `db:"outcome_price" json:"outcome_price,omitempty" obj:"price"`

	Quantity int `db:"quantity" json:"quantity,omitempty" obj:"quantity"`

	CategoryID *string `db:"category_id" json:"category_id,omitempty" obj:"category_id"`
	BrandID    *string `db:"brand_id" json:"brand_id,omitempty" obj:"brand_id"`

	Status string  `db:"status" json:"status,omitempty"`
	Rating float32 `db:"rating" json:"rating,omitempty" obj:"rating"`

	MainImage *string `db:"main_image" json:"main_image,omitempty" obj:"main_image"`

	ImageFiles []ProductMediaFiles `json:"image_files,omitempty" obj:"image_files"`
	VideoFiles []ProductMediaFiles `json:"video_files,omitempty" obj:"video_files"`

	CreatedAt time.Time  `db:"created_at" json:"created_at" obj:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at" obj:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type ProductMediaFiles struct {
	ID        string `db:"id" json:"id,omitempty"`
	ProductID string `db:"product_id" json:"product_id,omitempty"`
	MediaFile string `db:"media_file" json:"media_file,omitempty"`
}

type GetProductAllLimits struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
