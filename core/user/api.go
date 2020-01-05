package user

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/pkg/email"
	"github.com/novaladip/geldstroom-api-go/pkg/entity"
	errorsresponse "github.com/novaladip/geldstroom-api-go/pkg/errors"
)

// RegisterHandler ...
func RegisterHandler(r *gin.Engine, service Service) {
	res := resource{service}

	userRoute := r.Group("user")
	{
		userRoute.POST("/register", res.create)
		userRoute.POST("/login", res.login)
		userRoute.GET("/verify/email/:token", res.verifyEmail)
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

	u, err := r.service.Create(entity.User{
		Id:              entity.GenerateID(),
		Email:           dto.Email,
		Password:        dto.Password,
		IsActive:        true,
		JoinDate:        time.Now(),
		LastActivity:    time.Now(),
		IsEmailVerified: false,
	})

	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			c.JSON(http.StatusBadRequest, errorsresponse.BadRequestResponse{
				ErrorCode: ErrDuplicateEmailCode,
				Message:   ErrDuplicateEmail.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
		return
	}

	defer func() {
		t, err := r.service.CreateEmailVerification(u.Id)
		if err != nil {
			fmt.Println(err)
			return
		}
		_ = email.SendEmailVerification(u.Email, t)

	}()

	c.JSON(http.StatusCreated, u.GetWithoutPassword())

}

func (r resource) login(c *gin.Context) {
	var dto CredentialsDto
	_ = c.ShouldBind(&dto)
	if validate := dto.validate(); !validate.IsValid {
		c.JSON(http.StatusBadRequest, errorsresponse.ValidationError(ErrValidationFailedCode, ErrValidationFailed, validate.Error))
		return
	}

	u, err := r.service.Login(dto)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, errorsresponse.BadRequestResponse{
				ErrorCode: ErrInvalidCredentialsCode,
				Message:   ErrInvalidCredentials.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
		return
	}

	if ok := u.ValidatePassword(dto.Password); !ok {
		c.JSON(http.StatusBadRequest, errorsresponse.BadRequestResponse{
			ErrorCode: ErrInvalidCredentialsCode,
			Message:   ErrInvalidCredentials.Error(),
		})
		return
	}

	if !u.IsActive {
		c.JSON(http.StatusBadGateway, errorsresponse.BadRequestResponse{
			ErrorCode: ErrUserIsNotActiveCode,
			Message:   ErrUserIsNotActive.Error(),
		})
		return
	}

	if !u.IsEmailVerified {
		c.JSON(http.StatusBadGateway, errorsresponse.BadRequestResponse{
			ErrorCode: ErrEmailIsNotVerifiedCode,
			Message:   ErrEmailIsNotVerified.Error(),
		})
		return
	}

	token, err := r.service.GenerateJWT(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
	}
	c.JSON(http.StatusOK, gin.H{"Bearer": token})
}

func (r resource) verifyEmail(c *gin.Context) {
	tParam := c.Param("token")
	t, err := r.service.FindOneToken(tParam)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorsresponse.NotFound("Token not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
		return
	}

	if t.IsExpired() {
		c.JSON(http.StatusBadRequest, errorsresponse.BadRequestResponse{
			ErrorCode: ErrEmailVerificationExpiredCode,
			Message:   ErrEmailVerificationExpired.Error(),
		})
		return
	}

	if t.IsClaimed {
		c.JSON(http.StatusBadRequest, errorsresponse.BadRequestResponse{
			ErrorCode: ErrEmailVerificationAlreadyClaimedCode,
			Message:   ErrEmailVerificationAlreadyClaimed.Error(),
		})
		return
	}

	if err = r.service.VerifyEmail(t.UserId, t.Id); err != nil {
		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
