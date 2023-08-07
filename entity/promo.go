package entity

import (
	"time"

	"gorm.io/gorm"
)

type Promo struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name"`
	MinFee      float64        `json:"min_fee"`
	Discount    float64        `json:"discount"`
	MaxDiscount float64        `json:"max_discount"`
	Quota       int            `json:"quota"`
	Limited     bool           `json:"limited"`
	ExpDate     time.Time      `json:"exp_date"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}
