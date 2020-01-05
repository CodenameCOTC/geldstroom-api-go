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
	CreateEmailVerification(id string) (string, error)
	FindOneToken(token string) (entity.EmailVerification, error)
	VerifyEmail(userId, tokenId string) error
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

func (r repository) CreateEmailVerification(id string) (string, error) {
	e := entity.NewEmailVerification(id)
	stmt := `INSERT INTO token (id, token, expireAt, isClaimed, userId) VALUES(?, ?, ?, FALSE, ?) `
	_, err := r.DB.Exec(stmt, e.Id, e.Token, e.ExpireAt, e.UserId)
	if err != nil {
		return "", err
	}

	return e.Token, nil
}

func (r repository) FindOneToken(token string) (entity.EmailVerification, error) {
	ev := entity.EmailVerification{}
	stmt := `SELECT * FROM token WHERE token = ?`
	row := r.DB.QueryRow(stmt, token)
	err := row.Scan(&ev.Id, &ev.Token, &ev.ExpireAt, &ev.IsClaimed, &ev.UserId)
	if err != nil {
		return ev, err
	}
	return ev, nil
}

func (r repository) VerifyEmail(userId, tokenId string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	stmt := `UPDATE USER SET isEmailVerified = TRUE where id = ?`
	_, err = tx.Exec(stmt, userId)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	stmt = `UPDATE token SET isClaimed = TRUE WHERE id = ?`
	_, err = tx.Exec(stmt, tokenId)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
