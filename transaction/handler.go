package transaction

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/auth"
)

func GetTransactions(c *gin.Context) {
	user, ok := c.MustGet("JwtPayload").(auth.JwtPayload)

	if !ok {
		c.JSON(http.StatusUnauthorized, user)
		return
	}

	c.JSON(http.StatusOK, user)

}
