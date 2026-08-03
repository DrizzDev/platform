package sentry

import (
	sdk "github.com/getsentry/sentry-go"

	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/platform/configuration/reporting/sentry"
)

type Options struct {
	Build     build.Info
	Transport sdk.Transport
	Settings  sentry.Settings
}
