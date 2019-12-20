package auth

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/novaladip/geldstroom-api-go/logger"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Insert(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		logger.ErrorLog.Println(err)
		return err
	}

	stmt := `INSERT INTO user (email, password, isActive, joinDate, lastActivity) VALUES(?, ?, TRUE, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	_, err = h.Db.Exec(stmt, email, string(hashedPassword))

	if err != nil {
		logger.ErrorLog.Println(err)
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (h *Handler) Authenticate(credentials Credentials) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, password FROM user where email = ? AND isActive = TRUE"
	row := h.Db.QueryRow(stmt, credentials.Email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		logger.ErrorLog.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(credentials.Password))
	if err != nil {
		logger.ErrorLog.Println(err)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (h *Handler) Get(id int) (*UserModel, error) {
	u := &UserModel{}

	stmt := `SELECT id, email, isActive, joinDate, lastActivity, isEmailVerified FROM user where id = ?`

	err := h.Db.QueryRow(stmt, id).Scan(&u.ID, &u.Email, &u.IsActive, &u.JoinDate, &u.LastActivity, &u.IsEmailVerified)
	if err != nil {
		logger.ErrorLog.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return u, nil
}
