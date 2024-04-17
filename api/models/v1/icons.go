package models_v1

import "mime/multipart"

type AddIconToList struct {
	Name string                `form:"name" binding:"required"`
	Icon *multipart.FileHeader `form:"icon" binding:"required" swaggerignore:"true"`
}
