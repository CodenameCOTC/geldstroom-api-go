package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/helper"
)

type Handler struct {
	Db *sql.DB
}

func (h *Handler) Login(c *gin.Context) {
	var credentials Credentials

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

	err := h.Insert(credentials.Email, credentials.Password)

	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			c.JSON(http.StatusBadRequest, ErrDuplicateEmailDto)
			return
		}
		helper.ServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Register success, please check your email to verify your email.",
	})

}
