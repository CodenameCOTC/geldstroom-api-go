package user

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/novaladip/geldstroom-api-go/pkg/config"
	"github.com/novaladip/geldstroom-api-go/pkg/entity"
)

type Service interface {
	Create(entity.User) (entity.User, error)
	Login(dto CredentialsDto) (entity.User, error)
	GenerateJWT(user entity.User) (string, error)
	CreateEmailVerification(id string) (string, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) Create(u entity.User) (entity.User, error) {
	return s.repo.Create(u)

}

func (s service) Login(dto CredentialsDto) (entity.User, error) {
	return s.repo.FindOneByEmail(dto.Email)

}

func (s service) GenerateJWT(user entity.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 240).Unix(),
	}).SignedString([]byte(config.ConfigKey.SECRET))
}

func (s service) CreateEmailVerification(id string) (string, error) {
	return s.repo.CreateEmailVerification(id)
}
