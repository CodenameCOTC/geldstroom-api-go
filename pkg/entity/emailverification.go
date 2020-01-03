package entity

import "time"

type EmailVerification struct {
	Id        string    `json:"id"`
	Token     string    `json:"token"`
	ExpireAt  time.Time `json:"expireAt"`
	IsClaimed bool      `json:"isClaimed"`
	UserId    string    `json:"userId"`
}

func (ev EmailVerification) IsExpired() bool {
	return ev.ExpireAt.Before(time.Now())
}
