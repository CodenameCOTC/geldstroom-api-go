package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id              string    `json:"id"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	IsActive        bool      `json:"IsActive"`
	JoinDate        time.Time `json:"joinDate"`
	LastActivity    time.Time `json:"lastActivity"`
	IsEmailVerified bool      `json:"isEmailVerified"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u User) GetWithoutPassword() User {
	u.Password = ""
	return u
}
