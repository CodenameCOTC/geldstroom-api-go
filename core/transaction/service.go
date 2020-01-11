package transaction

import (
	"github.com/novaladip/geldstroom-api-go/pkg/entity"
)

type Service interface {
	Create(t entity.Transaction) (entity.Transaction, error)
	FindOneById(id, userId string) (entity.Transaction, error)
	DeleteOneById(id, userId string) error
	UpdateOneById(id, userId string, dto UpdateDto) (entity.Transaction, error)
	Get(p GetParam) ([]entity.Transaction, int, error)
	GetTotal(p GetTotalParam) (entity.TotalAmount, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) GetTotal(p GetTotalParam) (entity.TotalAmount, error) {
	return s.repo.GetTotal(p)
}

func (s service) Get(p GetParam) ([]entity.Transaction, int, error) {
	return s.repo.Get(p)
}

func (s service) Create(t entity.Transaction) (entity.Transaction, error) {
	return s.repo.Create(t)
}

func (s service) FindOneById(id, userId string) (entity.Transaction, error) {
	return s.repo.FindOneById(id, userId)
}

func (s service) DeleteOneById(id, userId string) error {
	return s.repo.DeleteOneById(id, userId)
}

func (s service) UpdateOneById(id, userId string, dto UpdateDto) (entity.Transaction, error) {
	return s.repo.UpdateOneById(id, userId, dto)
}
