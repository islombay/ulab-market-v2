package models

import "time"

type RoleModel struct {
	ID          string            `db:"id,primarykey" json:"id"`
	Name        string            `db:"name" json:"name"`
	Description *string           `db:"description" json:"description"`
	Permissions []PermissionModel `json:"permissions"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

func GetStringAddress(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func GetStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func GetIntAddress(i int) *int {
	return &i
}

type PermissionModel struct {
	ID          string  `db:"id,primarykey" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type AttachPermission struct {
	RoleID       string `db:"role_id" json:"role_id" binding:"required"`
	PermissionID string `db:"permission_id" json:"permission_id" binding:"required"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
