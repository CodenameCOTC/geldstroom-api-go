package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/novaladip/geldstroom-api-go/config"
)

type Claims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func (h *Handler) SignToken(id int, email string) (string, error) {
	key := config.GetKey()
	claims := &Claims{
		Id:    id,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().AddDate(0, 0, 7).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(key.SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
