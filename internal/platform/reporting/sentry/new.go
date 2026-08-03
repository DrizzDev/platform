package sentry

import (
	"log/slog"
	"sync"

	sdk "github.com/getsentry/sentry-go"
	slogsentry "github.com/samber/slog-sentry/v2"

	"github.com/DrizzDev/platform/internal/platform/logging"
)

func New(options Options) (Provider, error) {
	if !options.Settings.Enabled() {
		return Provider{}, nil
	}

	client, failure := sdk.NewClient(sdk.ClientOptions{
		Dsn:         options.Settings.DSN(),
		ServerName:  options.Build.Name(),
		Release:     options.Build.Version(),
		SampleRate:  options.Settings.Sample(),
		Environment: options.Settings.Environment(),

		AttachStacktrace: true,
		DisableLogs:      true,
		DisableMetrics:   true,
		EnableTracing:    false,
		MaxErrorDepth:    depth,
		MaxBreadcrumbs:   breadcrumbs,
		DataCollection:   options.privacy(),
		Integrations:     options.integrate,
		Transport:        options.Transport,
	})
	if failure != nil {
		return Provider{}, failure
	}

	hub := sdk.NewHub(client, sdk.NewScope())
	handler := slogsentry.Option{
		Hub:         hub,
		AddSource:   true,
		Level:       slog.LevelError,
		ReplaceAttr: logging.Policy{}.Handler(),
	}.NewSentryHandler()

	return Provider{handler: handler, client: client, once: &sync.Once{}, outcome: new(error)}, nil
}
