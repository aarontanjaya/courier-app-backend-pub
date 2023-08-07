package dto

import "time"

type PeriodRequest struct {
	StartDate *time.Time `form:"start_date,omitempty"`
	EndDate   *time.Time `form:"end_date,omitempty"`
}

type PaymentReport struct {
	Revenue       float64 `json:"revenue"`
	TotalDiscount float64 `json:"total_discount"`
	AvgSize       float64 `json:"average_size"`
	Count         float64 `json:"count"`
}

type ShippingReport struct {
	ShippingOrdered int `json:"shipping_count"`
}
