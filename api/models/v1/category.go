package models_v1

import "mime/multipart"

type CreateCategory struct {
	NameUz   string `json:"name_uz" binding:"required"`
	NameRu   string `json:"name_ru" binding:"required"`
	ParentID string `json:"parent_id"`
}

type ChangeCategoryImage struct {
	CategoryID string                `form:"category_id" binding:"required"`
	Image      *multipart.FileHeader `form:"image" binding:"required" swaggerignore:"true"`
}

type ChangeCategory struct {
	ID       string `json:"id" binding:"required"`
	NameUz   string `json:"name_uz" binding:"required"`
	NameRu   string `json:"name_ru" binding:"required"`
	ParentID string `json:"parent_id"`
}
