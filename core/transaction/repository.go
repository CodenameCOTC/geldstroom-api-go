package transaction

import (
	"database/sql"
	"errors"

	"github.com/novaladip/geldstroom-api-go/pkg/entity"
	"github.com/novaladip/geldstroom-api-go/pkg/getrange"
)

type Repository interface {
	Get(dateRange getrange.Range, userId string) ([]entity.Transaction, error)
	Create(t entity.Transaction) (entity.Transaction, error)
	FindOneById(id, userId string) (entity.Transaction, error)
	DeleteOneById(id, userId string) error
	UpdateOneById(id, userId string, dto UpdateDto) (entity.Transaction, error)
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

func (r repository) FindOneById(id, userId string) (entity.Transaction, error) {
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

func (r repository) DeleteOneById(id, userId string) error {
	stmt := `DELETE FROM transaction WHERE id = ? AND userId = ?`
	result, err := r.DB.Exec(stmt, id, userId)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrTransactionNotFound
	}

	return nil
}

func (r repository) UpdateOneById(id, userId string, dto UpdateDto) (entity.Transaction, error) {
	t := entity.Transaction{}
	stmt := `UPDATE transaction SET amount=?, category=?, type=?, description=? WHERE userId = ? AND id = ?`
	_, err := r.DB.Exec(stmt, dto.Amount, dto.Category, dto.Type, dto.Description, userId, id)
	if err != nil {
		return t, err
	}

	t, err = r.FindOneById(id, userId)
	if err != nil {
		return t, err
	}

	return t, nil
}

func (r repository) Get(dateRange getrange.Range, userId string) ([]entity.Transaction, error) {
	stmt := `SELECT * FROM transaction WHERE userId = ? AND createdAt BETWEEN ? AND ? ORDER BY updatedAt DESC LIMIT 10`
	rows, err := r.DB.Query(stmt, userId, dateRange.FirstDay, dateRange.LastDay)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	transactions := []entity.Transaction{}

	for rows.Next() {
		t := entity.Transaction{}
		err = rows.Scan(&t.Id, &t.Amount, &t.Description, &t.Category, &t.Type, &t.UserId, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil

}
