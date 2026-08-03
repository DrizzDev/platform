package logging_test

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/platform/configuration/logging"
	platform "github.com/DrizzDev/platform/internal/platform/logging"
)

func BenchmarkRecord(bench *testing.B) {
	logger := recorder{bench: bench}.logger()
	bench.ReportAllocs()
	bench.ResetTimer()
	for range bench.N {
		logger.Info(
			"operation completed",
			slog.String("operation", "MCP_SERVE"),
			slog.String("outcome", "SUCCESS"),
			slog.Duration("duration", time.Millisecond),
		)
	}
}

func BenchmarkRedaction(bench *testing.B) {
	logger := recorder{bench: bench}.logger()
	bench.ReportAllocs()
	bench.ResetTimer()
	for range bench.N {
		logger.Info(
			"request",
			slog.String("device", "pixel"),
			slog.String("access_token", "secret"),
		)
	}
}

type recorder struct {
	bench *testing.B
}

func (recorder recorder) logger() *slog.Logger {
	recorder.bench.Helper()
	settings, failure := logging.New("")
	if failure != nil {
		recorder.bench.Fatal(failure)
	}
	logger, failure := platform.New(platform.Options{
		Output:   io.Discard,
		Settings: settings,
		Build:    build.Read(),
	})
	if failure != nil {
		recorder.bench.Fatal(failure)
	}
	return logger
}
