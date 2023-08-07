package usecase_errors

import "errors"

var ErrEmailAlreadyExist = errors.New("email already exist")
var ErrReferralNotFound = errors.New("referral invalid")
