package dto

import "courier-app/entity"

type AddressResponsePayload struct {
	Records    *[]entity.Address `json:"records"`
	Search     string            `json:"search"`
	PageSize   int64             `json:"page_size"`
	TotalCount int64             `json:"total_count"`
	Page       int64             `json:"page"`
	PageCount  int64             `json:"page_count"`
}

type AddressUpdateRequest struct {
	ID             uint   `json:"id"`
	RecipientName  string `json:"recipient_name" binding:"required"`
	FullAddress    string `json:"full_address" binding:"required"`
	RecipientPhone string `json:"recipient_phone" binding:"required"`
	Label          string `json:"label" binding:"required"`
}

type AddressCreateRequest struct {
	RecipientName  string `json:"recipient_name" binding:"required"`
	FullAddress    string `json:"full_address" binding:"required"`
	RecipientPhone string `json:"recipient_phone" binding:"required"`
	Label          string `json:"label" binding:"required"`
}
