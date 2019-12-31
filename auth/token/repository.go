package token

import (
	"database/sql"
	"time"

	"github.com/novaladip/geldstroom-api-go/logger"
)

func Create(db *sql.DB, userId int) (string, error) {
	t := time.Now()

	stmt := `INSERT INTO token (token, expireAt, userId) VALUES(?, ?, ?)`

	token := generateRandomToken()

	_, err := db.Exec(stmt, token, t.Add(time.Hour*24), userId)
	if err != nil {
		logger.ErrorLog.Println(err)
		return "", err
	}

	return token, nil
}

func Renew(db *sql.DB, tokenId int) (string, error) {
	stmt := `UPDATE token SET token = ?, expireAt = ? WHERE id = ?`
	t := time.Now()
	token := generateRandomToken()

	_, err := db.Exec(stmt, token, t.Add(time.Hour*24), tokenId)
	if err != nil {
		logger.ErrorLog.Println(err)
		return "", err
	}
	return token, nil
}

func Get(db *sql.DB, token string) (*TokenModel, error) {
	t := &TokenModel{}
	stmt := `SELECT * FROM token WHERE token = ?`

	err := db.QueryRow(stmt, token).Scan(&t.Id, &t.Token, &t.ExpireAt, &t.UserId)

	if err != nil {
		return nil, err
	}

	ok := !t.ExpireAt.Before(time.Now())
	if !ok {
		return nil, ErrTokenExpired
	}

	return t, nil
}

func FindOneByUserIdAndRenew(db *sql.DB, id int) (string, error) {
	t := &TokenModel{}
	stmt := `SELECT * FROM token where userId = ?`

	err := db.QueryRow(stmt, id).Scan(&t.Id, &t.Token, &t.ExpireAt, &t.UserId)

	if err != nil {
		return "", err
	}

	ok := !t.ExpireAt.Before(time.Now())
	if !ok {
		newToken, err := Renew(db, t.Id)
		if err != nil {
			return "", err
		}
		return newToken, nil
	}

	return t.Token, nil
}
