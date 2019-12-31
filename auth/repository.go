package auth

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/novaladip/geldstroom-api-go/logger"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Insert(email, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		logger.ErrorLog.Println(err)
		return 0, err
	}

	stmt := `INSERT INTO user (email, password, isActive, joinDate, lastActivity) VALUES(?, ?, TRUE, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	result, err := h.Db.Exec(stmt, email, string(hashedPassword))

	if err != nil {
		logger.ErrorLog.Println(err)
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "email") {
				return 0, ErrDuplicateEmail
			}
		}
		return 0, err
	}

	lastId, _ := result.LastInsertId()

	return int(lastId), nil
}

func (h *Handler) Authenticate(credentials Credentials) (int, error) {
	var id int
	var hashedPassword []byte
	var isEmailVerified bool

	stmt := "SELECT id, password, isEmailVerified FROM user where email = ? AND isActive = TRUE"
	row := h.Db.QueryRow(stmt, credentials.Email)
	err := row.Scan(&id, &hashedPassword, &isEmailVerified)
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

	if !isEmailVerified {
		return 0, ErrEmailIsNotVerified
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

func (h *Handler) ValidateEmail(id int) error {
	u, err := h.Get(id)
	if err != nil {
		return err
	}

	if u.IsEmailVerified {
		return ErrEmailIsAlreadyVerified
	}

	stmt := `UPDATE user SET isEmailVerified = TRUE WHERE id = ?`
	_, err = h.Db.Exec(stmt, id)
	if err != nil {
		logger.InfoLog.Println(err)
		return nil
	}

	return nil
}

func (h *Handler) FindOneByEmail(email string) (*UserModel, error) {
	u := &UserModel{}
	stmt := `SELECT * FROM user WHERE email = ?`
	err := h.Db.QueryRow(stmt, email).Scan(&u.ID, &u.Email, &u.Password, &u.IsActive, &u.JoinDate, &u.LastActivity, &u.IsEmailVerified)
	if err != nil {
		logger.ErrorLog.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return u, nil
}
