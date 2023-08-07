package dto

type PaymentRequest struct {
	VoucherId     *uint
	PaymentId     uint
	TotalDiscount float64
	TotalCost     float64
	Amount        float64
	Description   string
}

type PaymentRequestBody struct {
	VoucherId uint `json:"voucher_id,omitempty"`
}

type BalanceResponse struct {
	UserId  uint
	Balance float64
}
