package middleware

import "database/sql"

type Guard struct {
	Db     *sql.DB
	Secret string
}

type AuthHeader struct {
	Authorization string
}
