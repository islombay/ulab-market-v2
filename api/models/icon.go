package models

import "time"

type IconModel struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	URL  string `db:"url" json:"url"`

	CreatedAt time.Time  `db:"created_at" json:"-"`
	UpdatedAt time.Time  `db:"updated_at" json:"-"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}
