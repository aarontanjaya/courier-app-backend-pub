package usecase_errors

import "errors"

var ErrUserNotAuthorized = errors.New("user not authorized")
var ErrRecordNotExist = errors.New("records doesn't exists")
var ErrInvalidReq = errors.New("invalid request")
