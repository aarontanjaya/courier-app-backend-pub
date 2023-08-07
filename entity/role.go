package entity

import (
	"gorm.io/gorm"
)

type Role struct {
	Name string
	gorm.Model
}
