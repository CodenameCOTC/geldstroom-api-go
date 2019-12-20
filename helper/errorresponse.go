package helper

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/logger"
)

type ErrorResponse struct {
	Message   string            `json:"message"`
	ErrorCode string            `json:"errorCode"`
	Error     map[string]string `json:"error"`
}

var InternalServerError = map[string]string{
	"message": "Internal Server Error",
}

var Unauthorized = map[string]string{
	"message": "Unauthorized",
}

func ServerError(c *gin.Context, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	logger.ErrorLog.Output(2, trace)
	c.JSON(http.StatusInternalServerError, InternalServerError)
}
