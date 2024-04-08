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
