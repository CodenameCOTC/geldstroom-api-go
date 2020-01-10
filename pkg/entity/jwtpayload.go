package entity

import "github.com/gin-gonic/gin"

type JwtPayload struct {
	Id    string
	Email string
}

func JwtPayloadFromRequest(c *gin.Context) JwtPayload {
	user, _ := c.MustGet("JwtPayload").(JwtPayload)
	return user
}
