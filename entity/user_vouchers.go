package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserVoucher struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	PromoId   uint           `json:"promo_id"`
	Promo     Promo          `json:"promo" gorm:"foreignKey:PromoId"`
	UserId    uint           `json:"user_id"`
	ExpDate   time.Time      `json:"exp_date"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
