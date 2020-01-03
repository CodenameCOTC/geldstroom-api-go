package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorsresponse "github.com/novaladip/geldstroom-api-go/pkg/errors"
)

// RegisterHandler ...
func RegisterHandler(r *gin.Engine, service Service) {
	res := resource{service}

	userRoute := r.Group("user")
	{
		userRoute.POST("/register", res.create)
		userRoute.POST("/login", res.login)
	}
}

type resource struct {
	service Service
}

func (r resource) create(c *gin.Context) {
	var dto CreateUserDto
	_ = c.ShouldBind(&dto)
	if validate := dto.validate(); !validate.IsValid {
		c.JSON(http.StatusBadRequest, errorsresponse.ValidationError(ErrValidationFailedCode, ErrValidationFailed, validate.Error))
		return
	}
	r.service.Create(c, dto)

}

func (r resource) login(c *gin.Context) {
	var dto CredentialsDto
	_ = c.ShouldBind(&dto)
	if validate := dto.validate(); !validate.IsValid {
		c.JSON(http.StatusBadRequest, errorsresponse.ValidationError(ErrValidationFailedCode, ErrValidationFailed, validate.Error))
		return
	}
	r.service.Login(c, dto)
}
