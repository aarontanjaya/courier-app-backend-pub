package dto

type RegisterResponse struct {
	Email string `json:"email"`
	Id    uint   `json:"id"`
	Name  string `json:"name"`
}

type RegisterRequest struct {
	Email           string `form:"email" json:"email" binding:"required,email"`
	Name            string `form:"name" json:"name" binding:"required"`
	Password        string `form:"password" json:"password" binding:"required"`
	Phone           string `form:"phone" json:"phone" binding:"required"`
	ClaimedReferral string `form:"referral" json:"claimed_referral" `
}
