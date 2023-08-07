package usecase_errors

import "errors"

var ErrPromoDiscountNegative = errors.New("discount field must not be negative")
var ErrPromoMinFeeNegative = errors.New("min_fee field must not be negative")
var ErrMaxDiscountNegative = errors.New("max_discount must not be negative")
var ErrQuotaNegative = errors.New("quota must not be negative")
var ErrPromoDiscountMax = errors.New("discount must be lower than or equal 100%")
var ErrInsufficientGachaQuota = errors.New("gacha quota insufficient")
