package dto

import "mime/multipart"

type ProfileDetails struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Photo        string `json:"photo"`
	Role         string `json:"role"`
	ReferralCode string `json:"referral"`
}

type UpdateProfileTextRequest struct {
	Name  string `form:"name" json:"name" binding:"required"`
	Email string `form:"email" json:"email" binding:"required,email"`
	Phone string `form:"phone" json:"phone" binding:"required"`
}

type UpdateProfileRequest struct {
	Name  string                `form:"name" json:"name" binding:"required"`
	Email string                `form:"email" json:"email" binding:"required,email"`
	Phone string                `form:"phone" json:"phone" binding:"required"`
	Photo *multipart.FileHeader `form:"photo,omitempty" json:"photo,omitempty"`
}

type UpdateProfileContentRequest struct {
	Name  string `form:"name" json:"name" validate:"required"`
	Email string `form:"email" json:"email" validate:"required,email"`
	Phone string `form:"phone" json:"phone" validate:"required"`
}

type UserProfile struct {
	ID    uint
	Name  string
	Email string
	Phone string
	Photo []byte
}
