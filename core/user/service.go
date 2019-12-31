package user

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/core/entity"
	errorsresponse "github.com/novaladip/geldstroom-api-go/core/errors"
	"github.com/novaladip/geldstroom-api-go/core/validator"
)

type Service interface {
	Create(c *gin.Context, dto CreateUserDto)
}

type User struct {
	entity.User
}

type CreateUserDto struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (dto CreateUserDto) validate() validator.Validate {
	v := validator.Validate{
		IsValid: true,
		Error:   make(map[string]string),
	}

	if !validator.EmailRX.MatchString(dto.Email) {
		v.Error["email"] = "Invalid email address"
	}

	if strings.TrimSpace(dto.Email) == "" {
		v.Error["email"] = "Email is cannot be empty"
	}

	if len(strings.TrimSpace(dto.Password)) < 6 {
		v.Error["password"] = "Password length must be greater than 6 characters"
	}

	if strings.TrimSpace(dto.Password) == "" {
		v.Error["password"] = "Password is cannot be empty"
	}

	if len(v.Error) > 0 {
		v.IsValid = false
	}

	return v
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) Create(c *gin.Context, dto CreateUserDto) {
	if validate := dto.validate(); !validate.IsValid {
		c.JSON(http.StatusBadRequest, errorsresponse.ValidationError(ErrValidationFailedCode, ErrValidationFailed, validate.Error))
		return
	}

	id := entity.GenerateID()
	now := time.Now()

	user, error := s.repo.Create(entity.User{
		Id:              id,
		Email:           dto.Email,
		Password:        dto.Password,
		IsActive:        true,
		JoinDate:        now,
		LastActivity:    now,
		IsEmailVerified: false,
	})

	if error != nil {
		if errors.Is(error, ErrDuplicateEmail) {
			c.JSON(http.StatusBadRequest, errorsresponse.BadRequestResponse{
				ErrorCode: ErrDuplicateEmailCode,
				Message:   ErrDuplicateEmail.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
		return
	}

	c.JSON(http.StatusCreated, user.GetWithoutPassword())
	return
}
