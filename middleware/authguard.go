package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/auth"
)

func (adb *AuthDb) AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		var authHeader AuthHeader
		authentication := auth.Authhentication{Db: adb.Db, Secret: ""}

		c.BindHeader(&authHeader)

		if strings.TrimSpace(authHeader.Authorization) == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "token not found",
			})
		}

		tokenString := strings.Replace(authHeader.Authorization, "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, auth.ErrInvalidCredentials
			} else if method != jwt.SigningMethodHS512 {
				return nil, auth.ErrInvalidCredentials
			}

			return []byte("s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge"), nil
		})

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			c.AbortWithStatusJSON(401, auth.ErrInvalidCredentialsDto)
		}

		userId := int(claims["id"].(float64))

		u, err := authentication.Get(userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.AbortWithStatusJSON(401, auth.ErrInvalidCredentialsDto)
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
		}

		if !u.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, auth.ErrInactiveUserDto)
		}

		if !u.IsEmailVerified {
			c.AbortWithStatusJSON(http.StatusUnauthorized, auth.ErrEmailIsNotVerifiedDto)
		}

		c.Next()
	}
}
