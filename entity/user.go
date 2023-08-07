package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Name        string
	Email       string
	Password    string
	Phone       string
	UserDetail  UserDetail
	UserVoucher []UserVoucher
	RoleDetail  Role `gorm:"foreignKey:Role"`
	Role        uint
	Photo       []byte
	PhotoFormat string
	ID          uint `gorm:"primaryKey;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type UserDetail struct {
	UserID               uint `gorm:"primaryKey;"`
	ReferralCode         string
	ClaimedReferral      string
	TotalTransactions    float64
	GachaQuota           uint
	Balance              float64
	ReferralStatus       uint
	ReferralStatusDetail ReferralStatus `gorm:"foreignKey:ReferralStatus"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt
}
