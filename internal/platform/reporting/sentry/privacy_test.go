package sentry_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	sdk "github.com/getsentry/sentry-go"

	"github.com/DrizzDev/platform/internal/platform/build"
	configuration "github.com/DrizzDev/platform/internal/platform/configuration/reporting/sentry"
	reporting "github.com/DrizzDev/platform/internal/platform/reporting/sentry"
)

const canary = "canary_secret_value"

func TestMetadata(test *testing.T) {
	test.Parallel()

	transport := &capture{}
	logger := slog.New(fixture{test: test}.provider(transport).Handler())
	logger.ErrorContext(
		context.Background(),
		"command failed",
		slog.Any("error", errors.New("device unavailable")),
		slog.String("access_token", canary),
		slog.Group("metadata",
			slog.String("session_id", "session_123"),
			slog.String("api_key", canary),
		),
	)

	events := transport.Events()
	if len(events) != 1 {
		test.Fatalf("events = %d", len(events))
	}
	payload, failure := json.Marshal(events[0])
	if failure != nil {
		test.Fatal(failure)
	}
	if strings.Contains(string(payload), canary) {
		test.Fatalf("sensitive value reached Sentry: %s", payload)
	}
}

func TestEnvironment(test *testing.T) {
	test.Parallel()

	transport := &capture{}
	logger := slog.New(fixture{test: test}.provider(transport).Handler())
	logger.ErrorContext(context.Background(), "command failed", slog.Any("error", errors.New("device unavailable")))

	events := transport.Events()
	if len(events) != 1 {
		test.Fatalf("events = %d", len(events))
	}
	for _, key := range []string{"device", "os", "runtime"} {
		if events[0].Contexts[key] != nil {
			test.Fatalf("Sentry event carried %q machine context", key)
		}
	}
	if len(events[0].Modules) != 0 {
		test.Fatalf("Sentry event carried the module inventory: %v", events[0].Modules)
	}
}

func TestCode(test *testing.T) {
	test.Parallel()

	transport := &capture{}
	logger := slog.New(fixture{test: test}.provider(transport).Handler())
	logger.ErrorContext(context.Background(), "command.failed")

	events := transport.Events()
	if len(events) != 1 {
		test.Fatalf("events = %d", len(events))
	}
	payload, failure := json.Marshal(events[0])
	if failure != nil {
		test.Fatal(failure)
	}
	if strings.Contains(string(payload), canary) {
		test.Fatalf("a code-only report carried unexpected content: %s", payload)
	}
	if len(events[0].Exception) != 0 {
		test.Fatalf("a code-only report attached an exception: %+v", events[0].Exception)
	}
}

func TestChain(test *testing.T) {
	test.Parallel()

	transport := &capture{}
	logger := slog.New(fixture{test: test}.provider(transport).Handler())
	inner := errors.New("stdio transport closed")
	logger.ErrorContext(
		context.Background(),
		"command failed",
		slog.Any("error", fmt.Errorf("serve: %w", inner)),
	)

	events := transport.Events()
	if len(events) != 1 || len(events[0].Exception) < 2 {
		test.Fatalf("exceptions = %+v", events)
	}
	outer := events[0].Exception[len(events[0].Exception)-1]
	if outer.Value != "serve: stdio transport closed" || outer.Stacktrace == nil {
		test.Fatalf("outer exception = %+v", outer)
	}
	if events[0].Exception[0].Value != "stdio transport closed" {
		test.Fatalf("inner exception = %+v", events[0].Exception[0])
	}
}

func TestClose(test *testing.T) {
	test.Parallel()

	provider := fixture{test: test}.provider(&capture{})
	if failure := provider.Close(context.Background()); failure != nil {
		test.Fatalf("first close: %v", failure)
	}
	if failure := provider.Close(context.Background()); failure != nil {
		test.Fatalf("repeated close: %v", failure)
	}
}

func TestTimeout(test *testing.T) {
	test.Parallel()

	provider := fixture{test: test}.provider(nil)
	slog.New(provider.Handler()).ErrorContext(
		context.Background(),
		"command failed",
		slog.Any("error", errors.New("device unavailable")),
	)
	expired, cancel := context.WithDeadline(context.Background(), time.Unix(0, 0))
	cancel()
	if failure := provider.Close(expired); failure == nil {
		test.Fatal("a flush that cannot complete before the deadline must surface an error")
	}
	if failure := provider.Close(expired); failure == nil {
		test.Fatal("repeated close after a timeout must remain bounded and reported")
	}
}

type fixture struct {
	test *testing.T
}

func (fixture fixture) provider(transport sdk.Transport) reporting.Provider {
	fixture.test.Helper()
	settings, failure := configuration.New(configuration.Input{
		DSN:         "https://public@example.invalid/1",
		Environment: "test",
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	provider, failure := reporting.New(reporting.Options{
		Settings:  settings,
		Build:     build.Read(),
		Transport: transport,
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return provider
}
