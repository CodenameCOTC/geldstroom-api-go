package report

import "github.com/getsentry/sentry-go"

import "github.com/novaladip/geldstroom-api-go/pkg/config"

func ErrorWrapperWithSentry(err error) error {
	if config.ConfigKey.APP_MODE == config.APP_MODE {
		defer sentry.CaptureException(err)
	}
	return err
}
