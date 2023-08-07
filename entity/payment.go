package entity

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	Status            bool           `json:"status"`
	TotalCost         float64        `json:"total_cost"`
	VoucherId         *uint          `json:"voucher_id"`
	TransactionRecord Transaction    `json:"record"`
	TotalDiscount     float64        `json:"total_discount"`
	ID                uint           `json:"id"`
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `json:"-"`
}
