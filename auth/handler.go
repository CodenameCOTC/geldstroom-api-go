package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Db *sql.DB
}

func (h *Handler) Login(c *gin.Context) {
	var credentials Credentials

	if err := c.ShouldBind(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequestDto)
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	token, err := h.SignToken(id, credentials.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": "bearer " + token,
	})
}

func (h *Handler) Register(c *gin.Context) {
	var credentials Credentials

	if err := c.ShouldBind(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequestDto)
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Register success, please check your email to verify your email.",
	})

}
