package transaction

import (
	"database/sql"
	"errors"

	"github.com/novaladip/geldstroom-api-go/core/entity"
)

type Repository interface {
	Create(t entity.Transaction) (entity.Transaction, error)
	FindOneById(id string, userId string) (entity.Transaction, error)
}

type repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return repository{db}
}

func (r repository) Create(t entity.Transaction) (entity.Transaction, error) {
	stmt := `INSERT INTO transaction (id, amount, description, category, type, userId) VALUE(?, ?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(stmt, t.Id, t.Amount, t.Description, t.Category, t.Type, t.UserId)

	if err != nil {
		return t, err
	}

	return t, nil
}

func (r repository) FindOneById(id string, userId string) (entity.Transaction, error) {
	stmt := `SELECT * FROM transaction WHERE id = ? AND userId = ?`
	row := r.DB.QueryRow(stmt, id, userId)
	t := entity.Transaction{}

	err := row.Scan(
		&t.Id,
		&t.Amount,
		&t.Description,
		&t.Category,
		&t.Type,
		&t.UserId,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, ErrTransactionNotFound
		}
		return t, err
	}
	return t, nil
}
