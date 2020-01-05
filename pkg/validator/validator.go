package validator

import "regexp"

type Validate struct {
	Error   map[string]string
	IsValid bool
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func New() Validate {
	return Validate{
		Error:   make(map[string]string),
		IsValid: true,
	}
}
