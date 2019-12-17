package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (adb *AuthDb) SignToken(id int, email string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   time.Now().AddDate(0, 0, 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(adb.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
