package models

import (
	"time"
)

type Product struct {
	ID      string `db:"id" json:"id"`
	Articul string `db:"articul" json:"articul,omitempty"`
	NameUz  string `db:"name_uz" json:"name_uz,omitempty"`
	NameRu  string `db:"name_ru" json:"name_ru,omitempty"`

	DescriptionUz string `db:"description_uz" json:"description_uz,omitempty"`
	DescriptionRu string `db:"description_ru" json:"description_ru,omitempty"`

	IncomePrice  float32 `db:"income_price" json:"income_price,omitempty"`
	OutcomePrice float64 `db:"outcome_price" json:"outcome_price,omitempty"`

	Quantity int `db:"quantity" json:"quantity,omitempty"`

	CategoryID *string `db:"category_id" json:"category_id,omitempty"`
	BrandID    *string `db:"brand_id" json:"brand_id,omitempty"`

	Status string  `db:"status" json:"status,omitempty"`
	Rating float32 `db:"rating" json:"rating,omitempty"`

	MainImage *string `db:"main_image" json:"main_image,omitempty"`

	ImageFiles []ProductMediaFiles `json:"image_files,omitempty"`
	VideoFiles []ProductMediaFiles `json:"video_files,omitempty"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type ProductMediaFiles struct {
	ID        string `db:"id" json:"id,omitempty"`
	ProductID string `db:"product_id" json:"product_id,omitempty"`
	MediaFile string `db:"media_file" json:"media_file,omitempty"`
}
