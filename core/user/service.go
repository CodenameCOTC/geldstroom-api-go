package user

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/pkg/entity"
	errorsresponse "github.com/novaladip/geldstroom-api-go/pkg/errors"
	"github.com/novaladip/geldstroom-api-go/pkg/config"
	"github.com/novaladip/geldstroom-api-go/pkg/validator"
)

type Service interface {
	Create(c *gin.Context, dto CreateUserDto)
	Login(c *gin.Context, dto CredentialsDto)
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

type CredentialsDto struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (dto CredentialsDto) validate() validator.Validate {
	v := validator.New()

	if !validator.EmailRX.MatchString(dto.Email) {
		v.Error["email"] = "Invalid email address"
	}

	if strings.TrimSpace(dto.Email) == "" {
		v.Error["email"] = "Email is cannot be empty"
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

}

func (s service) Login(c *gin.Context, dto CredentialsDto) {
	user, err := s.repo.FindOneByEmail(dto.Email)

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

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	if !user.IsEmailVerified {
		c.JSON(http.StatusBadGateway, errorsresponse.BadRequestResponse{
			ErrorCode: ErrEmailIsNotVerifiedCode,
			Message:   ErrEmailIsNotVerified.Error(),
		})
		return
	}

	token, err := s.generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
	}

	c.JSON(http.StatusOK, gin.H{"Bearer": token})

}

func (s service) generateJWT(user entity.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 240).Unix(),
	}).SignedString([]byte(config.ConfigKey.SECRET))
}
