package auth

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func (adb *Authhentication) Insert(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	stmt := `INSERT INTO user (email, password, isActive, joinDate, lastActivity) VALUES(?, ?, TRUE, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	_, err = adb.Db.Exec(stmt, email, string(hashedPassword))

	if err != nil {
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

func (adb *Authhentication) Authenticate(credentials Credentials) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, password FROM user where email = ? AND isActive = TRUE"
	row := adb.Db.QueryRow(stmt, credentials.Email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(credentials.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (adb *Authhentication) Get(id int) (*UserModel, error) {
	u := &UserModel{}

	stmt := `SELECT id, email, isActive, joinDate, lastActivity, isEmailVerified FROM user where id = ?`

	err := adb.Db.QueryRow(stmt, id).Scan(&u.ID, &u.Email, &u.IsActive, &u.JoinDate, &u.LastActivity, &u.IsEmailVerified)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return u, nil
}
