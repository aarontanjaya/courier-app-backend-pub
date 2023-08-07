package usecase_errors

import "errors"

var ErrReferralInvalid = errors.New("referral code invalid")
var ErrReferralFailed = errors.New("failed to generate referral code")
var ErrGeneratePassword = errors.New("failed to generate password hash")
var ErrCredsInvalid = errors.New("email or password invalid")
