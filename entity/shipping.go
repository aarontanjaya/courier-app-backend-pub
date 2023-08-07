package entity

import (
	"time"

	"gorm.io/gorm"
)

type Shipping struct {
	SizeId         uint           `json:"size_id"`
	Size           Size           `json:"size"`
	AddOn          []AddOn        `json:"add_on" gorm:"many2many:add_on_shippings"`
	CategoryId     uint           `json:"category_id"`
	Category       Category       `json:"category"`
	AddressId      uint           `json:"address_id"`
	Address        Address        `json:"address"`
	PaymentId      uint           `json:"payment_id"`
	Payment        Payment        `json:"payment"`
	StatusId       uint           `json:"status_id"`
	ShippingStatus ShippingStatus `json:"status" gorm:"foreignKey:StatusId"`
	ReviewComment  string         `json:"comment"`
	ReviewRating   int            `json:"rating"`
	UserId         uint           `json:"user_id"`
	ID             uint           `json:"id"`
	CreatedAt      time.Time      `json:"date"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-"`
}
