package report

import "github.com/getsentry/sentry-go"

func ErrorWrapperWithSentry(err error) error {
	defer sentry.CaptureException(err)
	return err
}
