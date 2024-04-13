package models_v1

import "mime/multipart"

type CreateCategory struct {
	Name     string `json:"name" binding:"required"`
	ParentID string `json:"parent_id"`
}

type CategoryTranslation struct {
	CategoryID   string `json:"category_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	LanguageCode string `json:"language_code" binding:"required"`
}

type ChangeCategoryImage struct {
	CategoryID string                `form:"category_id" binding:"required"`
	Image      *multipart.FileHeader `form:"image" binding:"required" swaggerignore:"true"`
}

type ChangeCategory struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	ParentID string `json:"parent_id"`
}

type DeleteCategoryRequest struct {
	CategoryID string `form:"category_id" binding:"required"`
	Language   string `form:"language" binding:"required"`
}
