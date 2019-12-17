package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (adb *Authhentication) Login(c *gin.Context) {
	var credentials Credentials

	if err := c.ShouldBind(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequestDto)
		return
	}

	id, err := adb.Authenticate(credentials)

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

	token, err := adb.SignToken(id, credentials.Email)
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

func (adb *Authhentication) Register(c *gin.Context) {
	var credentials Credentials

	if err := c.ShouldBind(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequestDto)
		return
	}

	err := adb.Insert(credentials.Email, credentials.Password)

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
