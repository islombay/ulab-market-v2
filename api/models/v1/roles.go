package models_v1

type CreateNewRole struct {
	Name        string `json:"name" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description" binding:"required"`
}

type AttachRoleToPermission struct {
	RoleID       string `json:"role_id" binding:"required"`
	PermissionID string `json:"permission_id" binding:"required"`
}
