package user

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/novaladip/geldstroom-api-go/pkg/config"
	"github.com/novaladip/geldstroom-api-go/pkg/entity"
)

type Service interface {
	Create(entity.User) (entity.User, error)
	FindOneByEmail(email string) (entity.User, error)
	GenerateJWT(user entity.User) (string, error)
	CreateEmailVerification(id string) (string, error)
	FindOneToken(token string) (entity.EmailVerification, error)
	VerifyEmail(userId, tokenId string) error
	FindTokenByUserId(id string) (entity.EmailVerification, error)
	RenewToken(id string) (entity.EmailVerification, error)
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

func (s service) FindOneByEmail(email string) (entity.User, error) {
	return s.repo.FindOneByEmail(email)

}

func (s service) CreateEmailVerification(id string) (string, error) {
	return s.repo.CreateEmailVerification(id)
}

func (s service) FindOneToken(token string) (entity.EmailVerification, error) {
	return s.repo.FindOneToken(token)
}

func (s service) FindTokenByUserId(id string) (entity.EmailVerification, error) {
	return s.repo.FindTokenByUserId(id)
}

func (s service) RenewToken(id string) (entity.EmailVerification, error) {
	return s.repo.RenewToken(id)
}

func (s service) VerifyEmail(userId, tokenId string) error {
	return s.repo.VerifyEmail(userId, tokenId)
}

func (s service) GenerateJWT(user entity.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 240).Unix(),
	}).SignedString([]byte(config.ConfigKey.SECRET))
}
