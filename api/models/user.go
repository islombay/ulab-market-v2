package models

import (
	"database/sql"
	"time"
)

type Client struct {
	ID          string         `json:"id" obj:"id"`
	Name        string         `json:"name" obj:""`
	PhoneNumber sql.NullString `json:"phone_number"`
	Email       sql.NullString `json:"email"`
	Password    string         `json:"password"`
	Verified    bool           `json:"verified" obj:"verified"`

	CreatedAt time.Time    `json:"created_at" obj:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Staff struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	PhoneNumber sql.NullString `json:"phone_number"`
	Email       sql.NullString `json:"email"`
	Password    string         `json:"password"`
	RoleID      string         `json:"role_id"`

	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
