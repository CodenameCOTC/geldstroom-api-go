package user

import "errors"

var (
	ErrValidationFailed     = errors.New("Validation failed")
	ErrValidationFailedCode = "USER_0001"
	ErrDuplicateEmail       = errors.New("Email is already taken")
	ErrDuplicateEmailCode   = "USER_0002"
)
