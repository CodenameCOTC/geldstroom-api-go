package config

import "os"

type key struct {
	SENDGRID_KEY string
	DB_DSN       string
	ADDR         string
	SECRET       string
}

const (
	dbDsnKey    = "DB_DSN"
	sendgridKey = "SENDGRID_KEY"
	addrKey     = "ADDR"
	secretKey   = "SECRET"
)

func GetKey() *key {
	return &key{
		DB_DSN:       os.Getenv(dbDsnKey),
		SENDGRID_KEY: os.Getenv(sendgridKey),
		ADDR:         os.Getenv(addrKey),
		SECRET:       os.Getenv(secretKey),
	}
}
