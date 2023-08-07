package usecase_errors

import "errors"

var ErrMinFeeNotReached = errors.New("transaction cost insufficient to use voucher")
var ErrVoucherExpired = errors.New("voucher expired")
var ErrVoucherInvalid = errors.New("voucher invalid")
var ErrInsufficientBalance = errors.New("balance insufficient")
var ErrPaymentNotExist = errors.New("bill not exist")
var ErrPaymentAlreadyPaid = errors.New("bill already paid")
