package user

import (
	"strings"

	"github.com/novaladip/geldstroom-api-go/pkg/validator"
)

type CreateUserDto struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (dto CreateUserDto) validate() validator.Validate {
	v := validator.Validate{
		IsValid: true,
		Error:   make(map[string]string),
	}

	if !validator.EmailRX.MatchString(dto.Email) {
		v.Error["email"] = "Invalid email address"
	}

	if strings.TrimSpace(dto.Email) == "" {
		v.Error["email"] = "Email is cannot be empty"
	}

	if len(strings.TrimSpace(dto.Password)) < 6 {
		v.Error["password"] = "Password length must be greater than 6 characters"
	}

	if strings.TrimSpace(dto.Password) == "" {
		v.Error["password"] = "Password is cannot be empty"
	}

	if len(v.Error) > 0 {
		v.IsValid = false
	}

	return v
}

type CredentialsDto struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (dto CredentialsDto) validate() validator.Validate {
	v := validator.New()

	if !validator.EmailRX.MatchString(dto.Email) {
		v.Error["email"] = "Invalid email address"
	}

	if strings.TrimSpace(dto.Email) == "" {
		v.Error["email"] = "Email is cannot be empty"
	}

	if strings.TrimSpace(dto.Password) == "" {
		v.Error["password"] = "Password is cannot be empty"
	}

	if len(v.Error) > 0 {
		v.IsValid = false
	}

	return v
}
