package transaction

import "strconv"

import "github.com/novaladip/geldstroom-api-go/helper"

func (h *Handler) insert(dto insertDto, userId int) (*TransactionModel, error) {
	stmt := `INSERT INTO transaction (amount, description, category, type, userId) VALUE(?, ?, ?, ?, ?)`
	result, err := h.Db.Exec(stmt, dto.Amount, dto.Description, dto.Category, dto.Type, userId)

	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	transaction, err := h.get(lastId, userId)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (h *Handler) get(transactionId int64, userId int) (*TransactionModel, error) {
	stmt := `SELECT * FROM transaction WHERE id = ? AND userId = ?`
	row := h.Db.QueryRow(stmt, transactionId, userId)
	t := &TransactionModel{}

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
		return nil, err
	}

	return t, nil
}

func (h *Handler) getTransaction(userId *int, r *Range) ([]*TransactionModel, error) {
	stmt := `SELECT * FROM transaction WHERE userId = ? AND createdAt BETWEEN ? and ? ORDER BY createdAt DESC LIMIT 10`
	rows, err := h.Db.Query(stmt, userId, r.firstDay, r.lastDay)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	transactions := []*TransactionModel{}

	for rows.Next() {
		t := &TransactionModel{}
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

func (h *Handler) update(tId string, dto updateDto, userId int) (*TransactionModel, error) {
	stmt := `UPDATE transaction SET amount=?, category=?, type=?, description=? WHERE userId = ? AND id= ?`
	_, err := h.Db.Exec(stmt, dto.Amount, dto.Category, dto.Type, dto.Description, userId, tId)
	if err != nil {
		return nil, err
	}

	id, _ := strconv.ParseInt(tId, 10, 64)

	t, err := h.get(id, userId)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (h *Handler) delete(tId string, userId int) error {
	stmt := `DELETE FROM transaction where id = ? AND userId = ?`
	result, err := h.Db.Exec(stmt, tId, userId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return helper.ErrSqlNoRow
	}

	return nil
}
