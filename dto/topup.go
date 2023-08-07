package dto

type TopUpRequestBody struct {
	Amount *float64 `json:"amount" binding:"required"`
}
