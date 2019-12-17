package auth

import (
	"database/sql"
	"errors"
	"time"
)

type UserModel struct {
	ID           int
	Email        string
	Password     string
	isActive     bool
	joinDate     time.Time
	lastActivity time.Time
}

type Credentials struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password"  binding:"required"`
}

type AuthDb struct {
	Db     *sql.DB
	Secret string
}

type ErrorDto struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

var (
	ErrInvalidCredentials     = errors.New("Invalids credentials")
	ErrInvalidCredentialsCode = "AUTH_0001"
	ErrDuplicateEmail         = errors.New("Email is already taken")
	ErrDuplicateEmailCode     = "AUTH_0002"
	ErrEmailIsNotVerified     = errors.New("Email is not verified yet")
	ErrEmailIsNotVerifiedCode = "AUTH_0003"
	ErrBadRequest             = errors.New("Bad request")
	ErrBadRequestCode         = "AUTH_0004"
)

var ErrInvalidCredentialsDto = &ErrorDto{
	Message:   ErrInvalidCredentials.Error(),
	ErrorCode: ErrInvalidCredentialsCode,
}

var ErrDuplicateEmailDto = &ErrorDto{
	Message:   ErrDuplicateEmail.Error(),
	ErrorCode: ErrDuplicateEmailCode,
}

var ErrBadRequestDto = &ErrorDto{
	Message:   ErrBadRequest.Error(),
	ErrorCode: ErrBadRequestCode,
}
