package entity

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	ID          uint           `json:"id"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func (Category) TableName() string {
	return "categories"
}
