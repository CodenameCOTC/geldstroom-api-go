package entity

import (
	"crypto/rand"
	"fmt"
	"time"
)

type EmailVerification struct {
	Id        string    `json:"id"`
	Token     string    `json:"token"`
	ExpireAt  time.Time `json:"expireAt"`
	IsClaimed bool      `json:"isClaimed"`
	UserId    string    `json:"userId"`
}

func NewEmailVerification(userId string) EmailVerification {
	return EmailVerification{
		Id:        GenerateID(),
		Token:     generateRandomToken(),
		ExpireAt:  time.Now().Add(time.Hour * 24),
		IsClaimed: false,
		UserId:    userId,
	}
}

func (ev EmailVerification) IsExpired() bool {
	return ev.ExpireAt.Before(time.Now())
}

func generateRandomToken() string {
	b := make([]byte, 20)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
