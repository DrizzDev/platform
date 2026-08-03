package observability

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DrizzDev/platform/internal/platform/logging"
	"github.com/DrizzDev/platform/internal/platform/reporting"
	"github.com/DrizzDev/platform/internal/platform/telemetry"
)

func New(scope context.Context, options Options) (Provider, error) {
	reporter, failure := reporting.New(reporting.Options{
		Build:    options.Build,
		Settings: options.Settings.Reporting(),
	})
	if failure != nil {
		return Provider{}, failure
	}
	diagnostics, failure := logging.New(logging.Options{
		Build:    options.Build,
		Output:   options.Output,
		Settings: options.Settings.Logging(),
	})
	if failure != nil {
		return Provider{}, errors.Join(failure, reporter.Close(scope))
	}
	provider, failure := telemetry.New(scope, telemetry.Options{
		Build:    options.Build,
		Settings: options.Settings.Telemetry(),
	})
	if failure != nil {
		return Provider{}, errors.Join(failure, reporter.Close(scope))
	}
	sink := slog.New(slog.DiscardHandler)
	if handler := reporter.Handler(); handler != nil {
		sink = slog.New(handler)
	}
	return Provider{
		diagnostics: diagnostics,
		sink:        sink,
		external:    slog.New(slog.DiscardHandler),
		reporting:   reporter,
		telemetry:   provider,
	}, nil
}
