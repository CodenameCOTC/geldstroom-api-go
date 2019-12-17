package auth

import (
	"database/sql"
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

type Credentials struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password"  binding:"required"`
}

type Authhentication struct {
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
	ErrInactiveUser           = errors.New("User is already deadactivated his/her account")
	ErrInactiveUserCode       = "AUTH_0005"
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
