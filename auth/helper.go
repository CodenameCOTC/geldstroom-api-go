package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func (adb *Authhentication) SignToken(id int, email string) (string, error) {
	claims := &Claims{
		Id:    id,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().AddDate(0, 0, 7).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(adb.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
