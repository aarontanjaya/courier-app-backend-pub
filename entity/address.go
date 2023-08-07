package entity

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	RecipientName  string         `json:"recipient_name"`
	FullAddress    string         `json:"full_address"`
	RecipientPhone string         `json:"recipient_phone"`
	UserId         uint           `json:"user_id"`
	Label          string         `json:"label"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-"`
}

func (Address) TableName() string {
	return "addresses"
}
