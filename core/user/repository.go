package user

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/novaladip/geldstroom-api-go/pkg/entity"
)

type Repository interface {
	Create(entity.User) (entity.User, error)
	// Delete(id string) error
	FindOneByEmail(email string) (entity.User, error)
	FindOneById(id string) (entity.User, error)
	// VerifyEmail(id string) error
	// Deactivate(id string) error
}

type repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return repository{db}
}

func (r repository) Create(user entity.User) (entity.User, error) {

	err := user.HashPassword()
	if err != nil {
		return user, err
	}

	stmt := `INSERT INTO user (id, email, password, isActive, isEmailVerified, joinDate, lastActivity) VALUES(?, ?, ?, TRUE, FALSE, ?, ?)`

	_, err = r.DB.Exec(stmt, user.Id, user.Email, user.Password, user.JoinDate, user.LastActivity)

	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "email") {
				return user, ErrDuplicateEmail
			}
		}
		return user, err
	}

	return user.GetWithoutPassword(), nil
}

func (r repository) FindOneByEmail(email string) (entity.User, error) {
	var user entity.User
	stmt := `SELECT * FROM user WHERE email = ?`
	row := r.DB.QueryRow(stmt, email)
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.IsActive, &user.JoinDate, &user.LastActivity, &user.IsEmailVerified)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrInvalidCredentials
		}
		return user, err
	}

	return user, nil
}

// FindOneByID ...
func (r repository) FindOneById(id string) (entity.User, error) {
	var user entity.User
	stmt := `SELECT * FROM user where id = ?`
	row := r.DB.QueryRow(stmt, id)
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.IsActive, &user.JoinDate, &user.LastActivity, &user.IsEmailVerified)

	if err != nil {
		return user, err
	}

	return user, nil
}
