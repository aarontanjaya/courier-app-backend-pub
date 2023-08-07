package entity

import "gorm.io/gorm"

type ReferralStatus struct {
	Message string
	gorm.Model
}

func (ReferralStatus) TableName() string {
	return "referral_statuses"
}
