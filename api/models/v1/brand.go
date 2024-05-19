package models_v1

import "mime/multipart"

type CreateBrand struct {
	Name  string                `json:"name" form:"name" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required" swaggerignore:"true"`
}

type ChangeBrand struct {
	ID    string                `json:"id" form:"name" binding:"required"`
	Name  string                `json:"name" form:"name" binding:"required"`
	Image *multipart.FileHeader `form:"image"`
}
