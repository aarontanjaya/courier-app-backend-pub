package dto

import (
	"courier-app/entity"
	"time"
)

type PromoRequest struct {
	Name        string    `json:"name" binding:"required"`
	MinFee      float64   `json:"min_fee" binding:"required"`
	Discount    float64   `json:"discount" binding:"required"`
	MaxDiscount float64   `json:"max_discount" binding:"required"`
	Quota       int       `json:"quota" binding:"required"`
	Limited     *bool     `json:"limited" binding:"required"`
	ExpDate     time.Time `json:"exp_date" binding:"required"`
}

type PromoResponsePayload struct {
	Records    *[]entity.Promo `json:"records"`
	Search     string          `json:"search"`
	PageSize   int64           `json:"page_size"`
	TotalCount int64           `json:"total_count"`
	Page       int64           `json:"page"`
	PageCount  int64           `json:"page_count"`
}
