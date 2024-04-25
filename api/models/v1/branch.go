package models_v1

import "time"

type CreateBranch struct {
	Name      string    `json:"name" binding:"required"`
	OpenTime  time.Time `json:"open_time" binding:"required"`
	CloseTime time.Time `json:"close_time" binding:"required"`
}

type ChangeBranch struct {
	ID        string    `json:"id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	OpenTime  time.Time `json:"open_time" binding:"required"`
	CloseTime time.Time `json:"close_time" binding:"required"`
}
