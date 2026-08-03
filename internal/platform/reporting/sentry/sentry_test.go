package sentry_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	sdk "github.com/getsentry/sentry-go"
	tracing "go.opentelemetry.io/otel/trace"

	"github.com/DrizzDev/platform/internal/platform/build"
	configuration "github.com/DrizzDev/platform/internal/platform/configuration/reporting/sentry"
	reporting "github.com/DrizzDev/platform/internal/platform/reporting/sentry"
)

func TestDisabled(test *testing.T) {
	test.Parallel()

	settings, failure := configuration.New(configuration.Input{})
	if failure != nil {
		test.Fatal(failure)
	}
	provider, failure := reporting.New(reporting.Options{
		Settings: settings,
		Build:    build.Read(),
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if provider.Handler() != nil {
		test.Fatal("disabled Sentry provided a handler")
	}
}

func TestError(test *testing.T) {
	test.Parallel()

	settings, failure := configuration.New(configuration.Input{
		DSN:         "https://public@example.invalid/1",
		Environment: "test",
	})
	if failure != nil {
		test.Fatal(failure)
	}
	transport := &capture{}
	provider, failure := reporting.New(reporting.Options{
		Settings:  settings,
		Build:     build.Read(),
		Transport: transport,
	})
	if failure != nil {
		test.Fatal(failure)
	}
	logger := slog.New(provider.Handler())
	identity := tracing.NewSpanContext(tracing.SpanContextConfig{
		TraceID: tracing.TraceID{1},
		SpanID:  tracing.SpanID{1},
	})
	scope := tracing.ContextWithSpanContext(context.Background(), identity)
	logger.ErrorContext(
		scope,
		"command failed",
		slog.Any("error", errors.New("device unavailable")),
		slog.String("access_token", "sensitive"),
	)
	events := transport.Events()
	if len(events) != 1 || len(events[0].Exception) != 1 {
		test.Fatalf("events = %+v", events)
	}
	if events[0].Exception[0].Value != "device unavailable" {
		test.Fatalf("exception = %+v", events[0].Exception[0])
	}
	if events[0].Contexts["trace"] == nil {
		test.Fatalf("contexts = %+v", events[0].Contexts)
	}
	extra := events[0].Contexts["extra"]
	if extra["access_token"] != "[REDACTED]" {
		test.Fatalf("extra = %+v", extra)
	}
	if failure := provider.Close(context.Background()); failure != nil {
		test.Fatal(failure)
	}
}

func TestPolicy(test *testing.T) {
	test.Parallel()

	settings, failure := configuration.New(configuration.Input{
		DSN:         "https://public@example.invalid/1",
		Environment: "test",
		Sample:      "0.5",
	})
	if failure != nil {
		test.Fatal(failure)
	}
	transport := &capture{}
	provider, failure := reporting.New(reporting.Options{
		Settings:  settings,
		Build:     build.Read(),
		Transport: transport,
	})
	if failure != nil {
		test.Fatal(failure)
	}
	options := transport.options
	privacy := options.DataCollection
	if privacy == nil {
		test.Fatal("data collection is not configured")
	}
	checks := []struct {
		name  string
		valid bool
	}{
		{name: "sample rate", valid: options.SampleRate == 0.5},
		{name: "Sentry tracing disabled", valid: !options.EnableTracing},
		{name: "Sentry logs disabled", valid: options.DisableLogs},
		{name: "Sentry metrics disabled", valid: options.DisableMetrics},
		{name: "user data disabled", valid: privacy.UserInfo.IsSet && !privacy.UserInfo.Value},
		{name: "cookies disabled", valid: privacy.Cookies.Mode == sdk.CollectionOff},
		{name: "query parameters disabled", valid: privacy.QueryParams.Mode == sdk.CollectionOff},
		{name: "request headers disabled", valid: privacy.HTTPHeaders.Request.Mode == sdk.CollectionOff},
		{name: "response headers disabled", valid: privacy.HTTPHeaders.Response.Mode == sdk.CollectionOff},
		{name: "HTTP bodies disabled", valid: len(privacy.HTTPBodies) == 0},
	}
	for _, check := range checks {
		if !check.valid {
			test.Error(check.name)
		}
	}
	integrations := options.Integrations(nil)
	if len(integrations) != 1 || integrations[0].Name() != "OTel" {
		test.Fatalf("integrations = %+v", integrations)
	}
	if failure := provider.Close(context.Background()); failure != nil {
		test.Fatal(failure)
	}
}

func TestLevel(test *testing.T) {
	test.Parallel()

	settings, failure := configuration.New(configuration.Input{
		DSN:         "https://public@example.invalid/1",
		Environment: "test",
	})
	if failure != nil {
		test.Fatal(failure)
	}
	provider, failure := reporting.New(reporting.Options{
		Settings: settings,
		Build:    build.Read(),
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if provider.Handler().Enabled(context.Background(), slog.LevelInfo) {
		test.Fatal("Sentry accepted an info record")
	}
	if !provider.Handler().Enabled(context.Background(), slog.LevelError) {
		test.Fatal("Sentry rejected an error record")
	}
	if failure := provider.Close(context.Background()); failure != nil {
		test.Fatal(failure)
	}
}

type capture struct {
	sdk.MockTransport
	options sdk.ClientOptions
}

func (capture *capture) Configure(options sdk.ClientOptions) {
	capture.options = options
}
