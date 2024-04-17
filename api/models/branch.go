package models

import "time"

type BranchModel struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}
