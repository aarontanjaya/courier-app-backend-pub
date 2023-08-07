package entity

import (
	"time"

	"gorm.io/gorm"
)

type ShippingStatus struct {
	Name      string         `json:"name"`
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (ShippingStatus) TableName() string {
	return "shipping_statuses"
}
