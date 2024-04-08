package models

type RoleModel struct {
	ID          string            `db:"id,primarykey" json:"id"`
	Name        string            `db:"name" json:"name"`
	Description *string           `db:"description" json:"description"`
	Permissions []PermissionModel `json:"permissions"`
}

func GetStringAddress(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func GetIntAddress(i int) *int {
	return &i
}

type PermissionModel struct {
	ID          string  `db:"id,primarykey" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description"`
}

type AttachPermission struct {
	RoleID       string `db:"role_id" json:"role_id" binding:"required"`
	PermissionID string `db:"permission_id" json:"permission_id" binding:"required"`
}
