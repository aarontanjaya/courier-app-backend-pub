package dto

import "courier-app/entity"

type ShippingResponsePayload struct {
	Records    *[]entity.Shipping `json:"records"`
	PageSize   int64              `json:"page_size"`
	TotalCount int64              `json:"total_count"`
	Page       int64              `json:"page"`
	PageCount  int64              `json:"page_count"`
}
type ShippingCreateRequest struct {
	CategoryId uint   `json:"category_id" binding:"required"`
	AddressId  uint   `json:"address_id" binding:"required"`
	SizeId     uint   `json:"size_id" binding:"required"`
	AddOnIds   []uint `json:"add_ons"`
}

type ShippingTableRequest struct {
	PageSize    int64   `form:"page_size"`
	Page        int64   `form:"page"`
	SizeIds     []int64 `form:"size_ids"`
	CategoryIds []int64 `form:"category_ids"`
	StatusIds   []int64 `form:"status_ids"`
}

type ShippingStatusUpdateRequest struct {
	StatusId uint `json:"status_id" binding:"required"`
}

type ShippingReviewRequest struct {
	ReviewComment string `json:"comment" binding:"required"`
	ReviewRating  uint   `json:"rating" binding:"required"`
	ShippingId    uint   `json:"shipping_id"`
}
