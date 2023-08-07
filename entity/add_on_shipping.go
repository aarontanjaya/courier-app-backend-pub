package entity

type AddOnShipping struct {
	ShippingId int `gorm:"primaryKey"`
	AddOnId    int `gorm:"primaryKey"`
}
