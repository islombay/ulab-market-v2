package models_v1

type UUIDResponse struct {
	ID string `json:"id" example:"438b3f0a-126b-4085-a8e9-525dfe0941e5"`
}

type ChangeAdminRequest struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone_number"`
	RoleID   string `json:"role_id"`
	Password string `json:"password"`
}
