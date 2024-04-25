package models

import "time"

type Category struct {
	ID       string      `json:"id" db:"id"`
	NameUz   string      `json:"name_uz" db:"name_uz"`
	NameRu   string      `json:"name_ru" db:"name_ru"`
	Image    *string     `json:"image,omitempty" db:"image"`
	IconID   *string     `json:"icon_id,omitempty" db:"icon_id"`
	ParentID *string     `json:"parent_id,omitempty" db:"parent_id"`
	Sub      []*Category `json:"subcategories,omitempty"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CategorySwagger struct {
	ID       string                `json:"id" example:"a2d70daf-b4ac-4198-a6a0-999447483c18"`
	NameUz   string                `json:"name_uz" example:"electronic"`
	NameRu   string                `json:"name_ru" example:"электроника"`
	Image    *string               `json:"image,omitempty" example:"https://firebasestorage.googleapis.com/v0/b/ulab-market.appspot.com/o/test%2Fcategory%2Fa2d70daf-b4ac-4198-a6a0-999447483c18?alt=media&token=test%2Fcategory%2Fa2d70daf-b4ac-4198-a6a0-999447483c18"`
	Icon     *string               `json:"icon_id,omitempty"`
	ParentID *string               `json:"parent_id,omitempty" example:""`
	Sub      []*SubCategorySwagger `json:"subcategories,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type SubCategorySwagger struct {
	ID       string  `json:"id"`
	NameUz   string  `json:"name_uz" example:"electronic"`
	NameRu   string  `json:"name_ru" example:"электроника"`
	Image    *string `json:"image,omitempty"`
	Icon     *string `json:"icon_id,omitempty"`
	ParentID *string `json:"parent_id,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
