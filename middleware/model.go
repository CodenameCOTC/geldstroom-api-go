package middleware

import "database/sql"

type Guard struct {
	Db *sql.DB
}

type AuthHeader struct {
	Authorization string
}
