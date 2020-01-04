package user

import "errors"

var (
	ErrValidationFailed       = errors.New("Validation failed")
	ErrValidationFailedCode   = "USER_0001"
	ErrDuplicateEmail         = errors.New("Email is already taken")
	ErrDuplicateEmailCode     = "USER_0002"
	ErrInvalidCredentials     = errors.New("Invalid credentials")
	ErrInvalidCredentialsCode = "USER_0003"
	ErrEmailIsNotVerified     = errors.New("Please verify your email first")
	ErrEmailIsNotVerifiedCode = "USER_0004"
	ErrUserIsNotActive        = errors.New("User is already deadactivated her/him account")
	ErrUserIsNotActiveCode    = "USER_0005"
)
