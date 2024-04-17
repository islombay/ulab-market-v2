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

type ClientAddress struct {
	ID       string  `db:"id"`
	ClientID string  `db:"client_id"`
	Long     float64 `db:"long"`
	Lat      float64 `db:"lat"`
	Name     string  `db:"name"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}
