package token

import (
	"errors"
	"time"
)

type TokenModel struct {
	Id       int       `json:"id"`
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expireAt"`
	UserId   int       `json:"userId"`
}

var (
	ErrTokenExpired = errors.New("Token is expired")
)
