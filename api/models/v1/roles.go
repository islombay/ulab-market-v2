package models_v1

type CreateNewRole struct {
	Name        string `json:"name" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description" binding:"required"`
}
