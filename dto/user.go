package dto

type UserResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Password     string  `json:"-"`
	Phone        string  `json:"phone"`
	Role         uint    `json:"role"`
	RoleName     string  `json:"role_name"`
	Balance      float64 `json:"balance"`
	Photo        string  `json:"photo"`
	ReferralCode string  `json:"referral"`
	GachaQuota   uint    `json:"gacha_quota"`
}

type UserClaims struct {
	UserId uint   `json:"id"`
	Scope  string `json:"scope"`
}

type UserBalanceResponse struct {
	Balance float64 `json:"balance"`
	UserId  uint    `json:"user_id"`
}

type UserGachaQuotaResponse struct {
	GachaQuota float64 `json:"quota"`
	UserId     uint    `json:"user_id"`
}
