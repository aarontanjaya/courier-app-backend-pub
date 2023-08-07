package entity

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	Description string         `json:"description"`
	Amount      float64        `json:"amount"`
	UserId      uint           `json:"user_id"`
	PaymentId   *uint          `json:"payment_id"`
	ID          uint           `json:"id"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}
