package sentry_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/DrizzDev/platform/internal/platform/build"
	configuration "github.com/DrizzDev/platform/internal/platform/configuration/reporting/sentry"
	"github.com/DrizzDev/platform/internal/platform/logging"
	reporting "github.com/DrizzDev/platform/internal/platform/reporting/sentry"
)

func BenchmarkSilent(bench *testing.B) {
	measure := reporter{bench: bench}
	measure.run(slog.New(measure.base()))
}

func BenchmarkReporting(bench *testing.B) {
	measure := reporter{bench: bench}
	settings, failure := configuration.New(configuration.Input{
		DSN:         "https://public@example.invalid/1",
		Environment: "test",
	})
	if failure != nil {
		bench.Fatal(failure)
	}
	provider, failure := reporting.New(reporting.Options{
		Settings:  settings,
		Build:     build.Read(),
		Transport: &capture{},
	})
	if failure != nil {
		bench.Fatal(failure)
	}
	measure.run(slog.New(slog.NewMultiHandler(measure.base(), provider.Handler())))
}

type reporter struct {
	bench *testing.B
}

func (reporter reporter) base() slog.Handler {
	return slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
		Level:       slog.LevelError,
		ReplaceAttr: logging.Policy{}.Handler(),
	})
}

func (reporter reporter) run(logger *slog.Logger) {
	reporter.bench.ReportAllocs()
	reporter.bench.ResetTimer()
	for range reporter.bench.N {
		logger.ErrorContext(
			context.Background(),
			"command failed",
			slog.Any("error", errors.New("stdio transport closed")),
			slog.String("access_token", "secret"),
		)
	}
}
