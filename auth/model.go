package auth

import (
	"errors"
	"time"
)

type UserModel struct {
	ID              int       `json:"id"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	IsActive        bool      `json:"isActive"`
	JoinDate        time.Time `json:"joinDate"`
	LastActivity    time.Time `json:"lastActivity"`
	IsEmailVerified bool      `json:"isEmailVerified"`
}

type JwtPayload struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type Credentials struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type ResendEmailDto struct {
	Email string `form:"email"`
}

type ErrorDto struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

var (
	ErrInvalidCredentials           = errors.New("Invalids credentials")
	ErrInvalidCredentialsCode       = "AUTH_0001"
	ErrDuplicateEmail               = errors.New("Email is already taken")
	ErrDuplicateEmailCode           = "AUTH_0002"
	ErrEmailIsNotVerified           = errors.New("Email is not verified yet")
	ErrEmailIsNotVerifiedCode       = "AUTH_0003"
	ErrBadRequest                   = errors.New("Bad request")
	ErrBadRequestCode               = "AUTH_0004"
	ErrInactiveUser                 = errors.New("User is already deadactivated his/her account")
	ErrInactiveUserCode             = "AUTH_0005"
	ErrFormFieldError               = errors.New("You may type a wrong answer to some field")
	ErrFormFieldErrorCode           = "AUTH_0006"
	ErrEmailVerificationExpired     = errors.New("Verification token is expired")
	ErrEmailVerificationExpiredCode = "AUTH_0007"
	ErrEmailIsAlreadyVerified       = errors.New("Email is already verified")
	ErrEmailIsAlreadyVerifiedCode   = "AUTH_0008"
	ErrUserNotFound                 = errors.New("User is not found")
	ErrUserNotFoundCode             = "AUTH_0009"
)

var ErrInvalidCredentialsDto = &ErrorDto{
	Message:   ErrInvalidCredentials.Error(),
	ErrorCode: ErrInvalidCredentialsCode,
}

var ErrDuplicateEmailDto = &ErrorDto{
	Message:   ErrDuplicateEmail.Error(),
	ErrorCode: ErrDuplicateEmailCode,
}

var ErrEmailIsNotVerifiedDto = &ErrorDto{
	Message:   ErrEmailIsNotVerified.Error(),
	ErrorCode: ErrEmailIsNotVerifiedCode,
}

var ErrInactiveUserDto = &ErrorDto{
	Message:   ErrInactiveUser.Error(),
	ErrorCode: ErrInactiveUserCode,
}

var ErrBadRequestDto = &ErrorDto{
	Message:   ErrBadRequest.Error(),
	ErrorCode: ErrBadRequestCode,
}

var ErrEmailVerificationExpiredDto = &ErrorDto{
	Message:   ErrEmailVerificationExpired.Error(),
	ErrorCode: ErrEmailVerificationExpiredCode,
}

var ErrEmailIsAlreadyVerfiedDto = &ErrorDto{
	Message:   ErrEmailIsAlreadyVerified.Error(),
	ErrorCode: ErrEmailIsAlreadyVerifiedCode,
}

var ErrUserNotFoundDto = &ErrorDto{
	Message:   ErrUserNotFound.Error(),
	ErrorCode: ErrUserNotFoundCode,
}
