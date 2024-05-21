package models_v1

import (
	"time"
)

type ClientOutput struct {
	ID          string `json:"id" obj:"id"`
	Name        string `json:"name" obj:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Verified    bool   `json:"verified" obj:"verified"`

	CreatedAt time.Time `json:"created_at" obj:"created_at"`
}

type ClientUpdate struct {
	ID          string     `json:"id" binding:"required"`
	Name        *string    `json:"name"`
	Surname     *string    `json:"surname"`
	PhoneNumber *string    `json:"phone_number"`
	Email       *string    `json:"email"`
	Gender      *string    `json:"gender"`
	BirthDate   *time.Time `json:"birth_date"`
}
