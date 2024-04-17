package models

import "time"

type Brand struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`

	CreatedAt time.Time  `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
