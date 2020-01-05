package auth

import (
	"database/sql"

	"github.com/novaladip/geldstroom-api-go/pkg/entity"
	"github.com/novaladip/geldstroom-api-go/core/user"
)

type Repository interface {
	CheckIsUserExist(id string) error
}

type repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return repository{db}
}

func (r repository) CheckIsUserExist(id string) error {
	var u entity.User
	stmt := `SELECT * FROM user where id = ?`
	row := r.DB.QueryRow(stmt, id)
	err := row.Scan(&u.Id, &u.Email, &u.Password, &u.IsActive, &u.JoinDate, &u.LastActivity, &u.IsEmailVerified)

	if err != nil {
		return err
	}

	if !u.IsEmailVerified {
		return user.ErrEmailIsNotVerified
	}

	return nil
}
