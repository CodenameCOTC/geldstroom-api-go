package middleware

import "database/sql"

type AuthDb struct{ Db *sql.DB }

type AuthHeader struct {
	Authorization string
}
