package auth

import (
	"regexp"
	"strings"
)

var emailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type credentialsValidator struct {
	Credentials Credentials
	error       map[string]string
}

type validator struct {
	isValid bool
	error   map[string]string
}

func newCredentialsValidator(c *Credentials) credentialsValidator {
	return credentialsValidator{
		Credentials: *c,
		error:       make(map[string]string),
	}
}

func (cv *credentialsValidator) validate() bool {
	if !emailRX.MatchString(cv.Credentials.Email) {
		cv.error["email"] = "Invalid email address"
	}

	if strings.TrimSpace(cv.Credentials.Email) == "" {
		cv.error["email"] = "Email field is cannot be empty"
	}

	if strings.TrimSpace(cv.Credentials.Password) == "" {
		cv.error["password"] = "Password field is cannot be empty"
	}

	if len(cv.error) > 0 {
		return false
	}

	return true
}

func (dto *ResendEmailDto) validate() *validator {
	var v = &validator{
		isValid: true,
		error:   make(map[string]string),
	}

	if !emailRX.MatchString(dto.Email) {
		v.error["email"] = "Invalid email address"
	}

	if strings.TrimSpace(dto.Email) == "" {
		v.error["email"] = "Email field is cannot be empty"
	}

	if len(v.error) > 0 {
		v.isValid = false
	}

	return v
}
