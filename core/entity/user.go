package entity

import "time"

type User struct {
	Id              string    `json:"id"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	IsActive        bool      `json:"IsActive"`
	JoinDate        time.Time `json:"joinDate"`
	LastActivity    time.Time `json:"lastActivity"`
	IsEmailVerified bool      `json:"isEmailVerified"`
}

func (u *User) GetWithoutPassword() *User {
	u.Password = ""
	return u
}
