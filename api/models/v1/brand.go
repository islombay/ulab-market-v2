package models_v1

type CreateBrand struct {
	Name string `json:"name" binding:"required"`
}

type ChangeBrand struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
