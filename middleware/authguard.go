package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/auth"
	"github.com/novaladip/geldstroom-api-go/config"
	"github.com/novaladip/geldstroom-api-go/helper"
)

func (g *Guard) AuthGuard() gin.HandlerFunc {
	key := config.GetKey()
	return func(c *gin.Context) {
		var authHeader AuthHeader
		authentication := auth.Handler{Db: g.Db}

		c.BindHeader(&authHeader)

		if strings.TrimSpace(authHeader.Authorization) == "" {
			c.AbortWithStatusJSON(401, &helper.Unauthorized)
			return
		}

		tokenString := strings.Replace(authHeader.Authorization, "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, auth.ErrInvalidCredentials
			} else if method != jwt.SigningMethodHS512 {
				return nil, auth.ErrInvalidCredentials
			}

			return []byte(key.SECRET), nil
		})

		if !token.Valid {
			c.AbortWithStatusJSON(401, &helper.Unauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.AbortWithStatusJSON(401, &helper.Unauthorized)
			return
		}

		userId := int(claims["id"].(float64))

		u, err := authentication.Get(userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.AbortWithStatusJSON(401, &helper.Unauthorized)
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, &helper.InternalServerError)
			return
		}

		c.Set("JwtPayload", auth.JwtPayload{
			Id:    u.ID,
			Email: u.Email,
		})

		if !u.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &auth.ErrInactiveUserDto)
			return
		}

		if !u.IsEmailVerified {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &auth.ErrEmailIsNotVerifiedDto)
			return
		}

		c.Next()
	}
}
