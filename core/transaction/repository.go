package transaction

import (
	"database/sql"
	"errors"

	"github.com/novaladip/geldstroom-api-go/pkg/entity"
	"github.com/novaladip/geldstroom-api-go/pkg/errors/report"
	"github.com/novaladip/geldstroom-api-go/pkg/getrange"
)

type GetParam struct {
	DateRange getrange.Range
	Page      int
	PerPage   int
	UserId    string
	Category  string
	Type      string
}

type GetTotalParam struct {
	Category string
	UserId   string
	Range    getrange.Range
}

type Repository interface {
	Get(p GetParam) ([]entity.Transaction, int, error)
	GetTotal(p GetTotalParam) (entity.TotalAmount, error)
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

func (r repository) GetTotal(p GetTotalParam) (entity.TotalAmount, error) {
	var stmt string
	var row *sql.Row
	var err error
	var ta entity.TotalAmount

	if p.Category == "ALL" || p.Category == "" {
		stmt = `SELECT ( 
			SELECT IFNULL(SUM(amount), 0) FROM transaction WHERE userId = ? AND type = "INCOME" AND createdAt BETWEEN ? AND ?
			) AS income,
			(
			SELECT IFNULL(SUM(amount), 0) FROM transaction WHERE userId = ? AND type = "EXPENSE" AND createdAt BETWEEN ? AND ?
			) AS expense`
		row = r.DB.QueryRow(stmt, p.UserId, p.Range.FirstDay, p.Range.LastDay, p.UserId, p.Range.FirstDay, p.Range.LastDay)

	} else {
		stmt = `SELECT ( 
			SELECT IFNULL(SUM(amount), 0) FROM transaction WHERE userId = ? AND type = "INCOME" AND category = ? AND createdAt BETWEEN ? AND ?
			) AS income,
			(
			SELECT IFNULL(SUM(amount), 0) FROM transaction WHERE userId = ? AND type = "EXPENSE" AND category = ? AND createdAt BETWEEN ? AND ?
			) AS expense`
		row = r.DB.QueryRow(stmt, p.UserId, p.Category, p.Range.FirstDay, p.Range.LastDay, p.UserId, p.Category, p.Range.FirstDay, p.Range.LastDay)

	}

	err = row.Scan(&ta.Income, &ta.Expense)
	if err != nil {
		return ta, err
	}

	return ta, nil
}

func GetQuery(p GetParam, db *sql.DB) (*sql.Rows, *sql.Row, error) {
	var stmt string
	var rows *sql.Rows
	var row *sql.Row

	if p.Category != "ALL" && p.Type != "ALL" {
		stmt = `SELECT * FROM transaction WHERE userId = ? AND category = ? AND type = ? AND createdAt BETWEEN ? AND ? ORDER BY createdAt DESC LIMIT ?, ?`
		rows, err := db.Query(
			stmt,
			p.UserId,
			p.Category,
			p.Type,
			p.DateRange.FirstDay,
			p.DateRange.LastDay,
			(p.Page-1)*p.PerPage,
			p.PerPage,
		)

		if err != nil {
			return nil, nil, err
		}

		stmt = `SELECT COUNT(*) FROM transaction WHERE userId = ? AND category = ? AND type = ?`
		row = db.QueryRow(stmt, p.UserId, p.Category, p.Type)

		return rows, row, nil
	}

	if p.Category == "ALL" && p.Type == "ALL" {
		stmt = `SELECT * FROM transaction WHERE userId = ? AND createdAt BETWEEN ? AND ? ORDER BY createdAt DESC LIMIT ?, ?`
		rows, err := db.Query(
			stmt,
			p.UserId,
			p.DateRange.FirstDay,
			p.DateRange.LastDay,
			(p.Page-1)*p.PerPage,
			p.PerPage,
		)

		if err != nil {
			return nil, nil, err
		}

		stmt = `SELECT COUNT(*) FROM transaction WHERE userId = ? `
		row = db.QueryRow(stmt, p.UserId)

		return rows, row, nil
	}

	if p.Category == "ALL" && p.Type != "ALL" {
		stmt = `SELECT * FROM transaction WHERE userId = ? AND type = ? AND createdAt BETWEEN ? AND ? ORDER BY createdAt DESC LIMIT ?, ?`
		rows, err := db.Query(
			stmt,
			p.UserId,
			p.Type,
			p.DateRange.FirstDay,
			p.DateRange.LastDay,
			(p.Page-1)*p.PerPage,
			p.PerPage,
		)

		if err != nil {
			return nil, nil, err
		}

		stmt = `SELECT COUNT(*) FROM transaction WHERE userId = ? AND type = ?`
		row = db.QueryRow(stmt, p.UserId, p.Type)

		return rows, row, nil
	}

	if p.Category != "ALL" && p.Type == "ALL" {
		stmt = `SELECT * FROM transaction WHERE userId = ? AND category = ? AND createdAt BETWEEN ? AND ? ORDER BY createdAt DESC LIMIT ?, ?`
		rows, err := db.Query(
			stmt,
			p.UserId,
			p.Category,
			p.DateRange.FirstDay,
			p.DateRange.LastDay,
			(p.Page-1)*p.PerPage,
			p.PerPage,
		)

		if err != nil {
			return nil, nil, err
		}

		stmt = `SELECT COUNT(*) FROM transaction WHERE userId = ? AND category = ?`
		row = db.QueryRow(stmt, p.UserId, p.Category)

		return rows, row, nil
	}

	return rows, row, nil
}

func (r repository) Get(p GetParam) ([]entity.Transaction, int, error) {
	rows, row, err := GetQuery(p, r.DB)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	transactions := []entity.Transaction{}

	for rows.Next() {
		t := entity.Transaction{}
		err = rows.Scan(&t.Id, &t.Amount, &t.Description, &t.Category, &t.Type, &t.CreatedAt, &t.UpdatedAt, &t.UserId)
		if err != nil {
			return nil, 0, report.ErrorWrapperWithSentry(err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, report.ErrorWrapperWithSentry(err)
	}

	var count int
	if err = row.Scan(&count); err != nil {
		return nil, 0, report.ErrorWrapperWithSentry(err)
	}

	return transactions, count, nil
}

func (r repository) Create(t entity.Transaction) (entity.Transaction, error) {
	stmt := `INSERT INTO transaction (id, amount, description, category, type, userId) VALUE(?, ?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(stmt, t.Id, t.Amount, t.Description, t.Category, t.Type, t.UserId)

	if err != nil {
		return t, report.ErrorWrapperWithSentry(err)
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
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.UserId,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, report.ErrorWrapperWithSentry(ErrTransactionNotFound)
		}
		return t, report.ErrorWrapperWithSentry(err)
	}
	return t, nil
}

func (r repository) DeleteOneById(id, userId string) error {
	stmt := `DELETE FROM transaction WHERE id = ? AND userId = ?`
	result, err := r.DB.Exec(stmt, id, userId)

	if err != nil {
		return report.ErrorWrapperWithSentry(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return report.ErrorWrapperWithSentry(err)
	}

	if affected == 0 {
		return report.ErrorWrapperWithSentry(ErrTransactionNotFound)
	}

	return nil
}

func (r repository) UpdateOneById(id, userId string, dto UpdateDto) (entity.Transaction, error) {
	t := entity.Transaction{}
	stmt := `UPDATE transaction SET amount=?, category=?, type=?, description=? WHERE userId = ? AND id = ?`
	_, err := r.DB.Exec(stmt, dto.Amount, dto.Category, dto.Type, dto.Description, userId, id)
	if err != nil {
		return t, report.ErrorWrapperWithSentry(err)
	}

	t, err = r.FindOneById(id, userId)
	if err != nil {
		return t, report.ErrorWrapperWithSentry(err)
	}

	return t, nil
}
