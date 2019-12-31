package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/auth/token"
	"github.com/novaladip/geldstroom-api-go/helper"
	sendmail "github.com/novaladip/geldstroom-api-go/mail"
)

type Handler struct {
	Db *sql.DB
}

func (h *Handler) Login(c *gin.Context) {
	var credentials Credentials
	c.ShouldBind(&credentials)

	cv := newCredentialsValidator(&credentials)

	ok := cv.validate()
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":   ErrFormFieldError.Error(),
			"errorCode": ErrFormFieldErrorCode,
			"error":     cv.error,
		})
		return
	}

	id, err := h.Authenticate(credentials)

	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, ErrInvalidCredentialsDto)
			return
		}
		if errors.Is(err, ErrEmailIsNotVerified) {
			c.JSON(http.StatusBadRequest, ErrEmailIsNotVerifiedDto)
			return
		}
		helper.ServerError(c, err)
		return
	}

	token, err := h.SignToken(id, credentials.Email)
	if err != nil {
		helper.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": "Bearer " + token,
	})

}

func (h *Handler) Register(c *gin.Context) {
	var credentials Credentials

	c.ShouldBind(&credentials)

	cv := newCredentialsValidator(&credentials)
	ok := cv.validate()
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":   ErrFormFieldError.Error(),
			"errorCode": ErrFormFieldErrorCode,
			"error":     cv.error,
		})
		return
	}

	id, err := h.Insert(credentials.Email, credentials.Password)

	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			c.JSON(http.StatusBadRequest, ErrDuplicateEmailDto)
			return
		}
		helper.ServerError(c, err)
		return
	}

	t, _ := token.Create(h.Db, id)

	defer sendmail.Send(credentials.Email, "<p>Verify your email by clicking<a href=''> Here</a><p>"+t)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Register success, please check your email to verify your email.",
	})
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	paramToken := c.Param("token")

	t, err := token.Get(h.Db, paramToken)
	if err != nil {
		if errors.Is(err, token.ErrTokenExpired) {
			c.JSON(http.StatusBadRequest, ErrEmailVerificationExpiredDto)
			return
		}
		helper.ServerError(c, err)
		return
	}

	err = h.ValidateEmail(t.UserId)
	if err != nil {
		if errors.Is(err, ErrEmailIsAlreadyVerified) {
			c.JSON(http.StatusBadRequest, ErrEmailIsAlreadyVerfiedDto)
			return
		}
		helper.ServerError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Email verified, now you can login using your email address",
	})

}

func (h *Handler) ResendEmailVerification(c *gin.Context) {
	var dto ResendEmailDto

	c.ShouldBind(&dto)

	v := dto.validate()

	if !v.isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed",
			"error":   v.error,
		})
		return
	}

	user, err := h.FindOneByEmail(dto.Email)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusNotFound, ErrUserNotFoundDto)
			return
		}
		helper.ServerError(c, err)
		return
	}

	if user.IsEmailVerified {
		c.JSON(http.StatusBadRequest, ErrEmailIsAlreadyVerfiedDto)
		return
	}

	t, err := token.FindOneByUserIdAndRenew(h.Db, user.ID)
	if err != nil {
		helper.ServerError(c, err)
		return
	}

	defer sendmail.Send(user.Email, "<p>Verify your email by clicking<a href=''> Here</a><p>"+t)

	c.JSON(http.StatusNoContent, gin.H{})
}
