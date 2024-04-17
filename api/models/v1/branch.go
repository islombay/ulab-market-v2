package models_v1

type CreateBranch struct {
	Name string `json:"name" binding:"required"`
}

type ChangeBranch struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
