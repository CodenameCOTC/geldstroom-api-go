package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/pkg/config"
	"github.com/novaladip/geldstroom-api-go/pkg/entity"
	errorsresponse "github.com/novaladip/geldstroom-api-go/pkg/errors"
	"github.com/novaladip/geldstroom-api-go/pkg/errors/report"
)

type Middleware interface {
	AuthGuard() gin.HandlerFunc
}

type authMiddleware struct {
	repo Repository
}

func NewMiddleware(repo Repository) Middleware {
	return authMiddleware{repo}
}

type authHeader struct {
	Authorization string
}

func (am authMiddleware) AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := config.ConfigKey.SECRET
		var header authHeader

		_ = c.BindHeader(&header)

		if (strings.TrimSpace(header.Authorization)) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorsresponse.Unauthorized(""))
			return
		}

		tokenString := strings.Replace(header.Authorization, "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("")
			} else if method != jwt.SigningMethodHS256 {
				return nil, errors.New("")
			}

			return []byte(secretKey), nil
		})

		if err != nil {
			_ = report.ErrorWrapperWithSentry(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorsresponse.Unauthorized(""))
			return
		}

		if !token.Valid {
			_ = report.ErrorWrapperWithSentry(errors.New("Invalid bearer token"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorsresponse.Unauthorized(""))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorsresponse.Unauthorized(""))
			return
		}

		userId := fmt.Sprintf("%v", claims["id"])

		if err = am.repo.CheckIsUserExist(userId); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorsresponse.Unauthorized(""))
			return
		}

		c.Set("JwtPayload", entity.JwtPayload{
			Id:    fmt.Sprintf("%v", claims["id"]),
			Email: fmt.Sprintf("%v", claims["email"]),
		})

		c.Next()
	}
}
