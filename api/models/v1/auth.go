package models_v1

type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"Islombek"`
	Email    string `json:"email" binding:"required" example:"example@example.com"`
	Phone    string `json:"phone_number" binding:"required" example:"998912345678"`
	Password string `json:"password" binding:"required" example:"password_1"`
}

type RequestCode struct {
	Email    string `json:"email,omitempty" `
	Phone    string `json:"phone_number,omitempty"`
	NeedCode bool   `json:"need_code"`
}

type VerifyCodeRequest struct {
	Source string `json:"source" binding:"required" example:"example@example.com"`
	Type   string `json:"type" binding:"required" example:"email" enums:"email,phone_number"`
	Code   string `json:"code" binding:"required" example:"111111"`
}

type LoginRequest struct {
	Source string `json:"source" binding:"required"  example:"998912345678"`
	Type   string `json:"type" binding:"required"  example:"phone_number"`
}

type LoginAdminRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Token struct {
	Token string `json:"token"`
}

type ResponseID struct {
	ID string `json:"id"`
}

type RequestCodeRequest struct {
	Source string `json:"source" binding:"required" example:"example@example.com"`
	Type   string `json:"type" binding:"required"  example:"email" enums:"email,phone_number"`
}

type ChangePassword struct {
	Source   string `json:"source" binding:"required"  example:"example@example.com"`
	Type     string `json:"type" binding:"required"  example:"email" enums:"email,phone_number"`
	Code     string `json:"code" binding:"required"  example:"111111"`
	Password string `json:"password" binding:"required"  example:"password_1"`
}
