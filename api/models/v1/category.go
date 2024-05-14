package models_v1

import "mime/multipart"

type CreateCategory struct {
	NameUz   string                `form:"name_uz" json:"name_uz" binding:"required"`
	NameRu   string                `form:"name_ru" json:"name_ru" binding:"required"`
	IconID   *string               `form:"icon_id" json:"icon_id"`
	Image    *multipart.FileHeader `form:"image" swaggerignore:"true"`
	ParentID string                `form:"parent_id" json:"parent_id"`
}

type ChangeCategoryImage struct {
	CategoryID string                `form:"category_id" binding:"required"`
	IconID     *string               `form:"icon_id"`
	Image      *multipart.FileHeader `form:"image" swaggerignore:"true"`
}

type ChangeCategory struct {
	ID       string `json:"id" binding:"required"`
	NameUz   string `json:"name_uz" binding:"required"`
	NameRu   string `json:"name_ru" binding:"required"`
	ParentID string `json:"parent_id"`
}

type GetAllCategory struct {
	OnlySub bool `form:"only_sub"`
}
