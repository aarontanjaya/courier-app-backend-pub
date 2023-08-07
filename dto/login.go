package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserId      uint   `json:"user_id"`
	AccessToken string `json:"-"`
	Scope       string `json:"scope"`
}
