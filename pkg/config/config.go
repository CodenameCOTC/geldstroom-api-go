package config

import "os"

type Key struct {
	SENDGRID_KEY string
	DB_DSN       string
	ADDR         string
	SECRET       string
	SENTRY_DSN   string
}

const (
	dbDsnKey    = "DB_DSN"
	sendgridKey = "SENDGRID_KEY"
	addrKey     = "ADDR"
	secretKey   = "SECRET"
	sentryKey   = "SENTRY_DSN"
)

var ConfigKey Key

func LoadKey() {
	ConfigKey = Key{
		DB_DSN:       os.Getenv(dbDsnKey),
		SENDGRID_KEY: os.Getenv(sendgridKey),
		ADDR:         os.Getenv(addrKey),
		SECRET:       os.Getenv(secretKey),
		SENTRY_DSN:   os.Getenv(sentryKey),
	}
}
