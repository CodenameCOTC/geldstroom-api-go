package transaction

import "github.com/novaladip/geldstroom-api-go/core/entity"

type Service interface {
	Create(t entity.Transaction) (entity.Transaction, error)
	FindOneById(id string, userId string) (entity.Transaction, error)
	DeleteOneById(id string, userId string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) Create(t entity.Transaction) (entity.Transaction, error) {
	return s.repo.Create(t)
}

func (s service) FindOneById(id string, userId string) (entity.Transaction, error) {
	return s.repo.FindOneById(id, userId)
}

func (s service) DeleteOneById(id string, userId string) error {
	return s.repo.DeleteOneById(id, userId)
}
