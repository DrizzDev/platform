package configuration_test

import (
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/platform/configuration"
	"github.com/DrizzDev/platform/internal/platform/configuration/logging"
	"github.com/DrizzDev/platform/internal/platform/configuration/telemetry"
)

func TestDefaults(test *testing.T) {
	test.Parallel()

	settings, failure := configuration.New(nil).Load()
	if failure != nil {
		test.Fatal(failure)
	}
	if settings.Logging().Level() != logging.Off {
		test.Fatalf("level = %v", settings.Logging().Level())
	}
	if settings.Telemetry().Exporter() != telemetry.None {
		test.Fatalf("exporter = %q", settings.Telemetry().Exporter())
	}
	if settings.Reporting().Sentry().Enabled() {
		test.Fatal("sentry enabled without a DSN")
	}
	if settings.Reporting().Sentry().Sample() != 1 {
		test.Fatalf("sample = %v", settings.Reporting().Sentry().Sample())
	}
}

func TestOTLP(test *testing.T) {
	test.Parallel()

	settings, failure := configuration.New([]string{
		"DRIZZ_LOG_LEVEL=debug",
		"DRIZZ_TELEMETRY_EXPORTER=otlp",
		"DRIZZ_TELEMETRY_ENDPOINT=http://127.0.0.1:4318",
	}).Load()
	if failure != nil {
		test.Fatal(failure)
	}
	if settings.Logging().Level() != logging.Debug {
		test.Fatalf("level = %v", settings.Logging().Level())
	}
	if settings.Telemetry().Exporter() != telemetry.OTLP {
		test.Fatalf("exporter = %q", settings.Telemetry().Exporter())
	}
}

func TestIdentity(test *testing.T) {
	test.Parallel()

	settings, failure := configuration.New([]string{
		"DRIZZ_AUTH0_ISSUER=https://drizz-dev.eu.auth0.com/",
		"DRIZZ_AUTH0_CLIENT=platform_cli",
		"DRIZZ_AUTH0_AUDIENCE=https://platform.drizz.dev",
		"DRIZZ_AUTH0_REDIRECT=http://127.0.0.1:8490/callback",
		"DRIZZ_AUTH0_SCOPES=openid,offline_access",
		"DRIZZ_SESSION=LOCAL",
	}).Load()
	if failure != nil {
		test.Fatal(failure)
	}
	identity := settings.Identity()
	if identity.Issuer() != "https://drizz-dev.eu.auth0.com/" || identity.Client() != "platform_cli" {
		test.Fatalf("identity = %+v", identity)
	}
	if identity.Audience() != "https://platform.drizz.dev" || identity.Session() != "LOCAL" {
		test.Fatalf("identity = %+v", identity)
	}
	if len(identity.Scopes()) != 2 || identity.Scopes()[0] != "openid" {
		test.Fatalf("scopes = %v", identity.Scopes())
	}
}

func TestAmbient(test *testing.T) {
	test.Parallel()

	settings, failure := configuration.New([]string{
		"SENTRY_DSN=https://public@inherited.example/1",
		"SENTRY_SAMPLE_RATE=0.1",
		"OTEL_SERVICE_NAME=other",
		"OTEL_EXPORTER_OTLP_ENDPOINT=https://collector.example",
	}).Load()
	if failure != nil {
		test.Fatalf("inherited third-party variables were not ignored: %v", failure)
	}
	if settings.Reporting().Sentry().Enabled() {
		test.Fatal("an inherited SENTRY_DSN activated Drizz reporting")
	}
	if settings.Telemetry().Exporter() != telemetry.None {
		test.Fatalf("an inherited OTEL variable changed telemetry: %q", settings.Telemetry().Exporter())
	}
}

func TestConfidential(test *testing.T) {
	test.Parallel()

	const secret = "confidential_dsn_token"
	cases := [][]string{
		{"DRIZZ_SENTRY_DSN=http://public@" + secret + ".example/1"},
		{"DRIZZ_SENTRY_SAMPLE_RATE=" + secret},
		{"DRIZZ_TELEMETRY_EXPORTER=OTLP", "DRIZZ_TELEMETRY_ENDPOINT=https://user:" + secret + "@example.com"},
		{"DRIZZ_LOG_LEVEL=" + secret},
		{"DRIZZ_TELEMETRY_EXPORTER=" + secret},
		{"DRIZZ_SENTRY_ENVIRONMENT=" + secret + "/../etc"},
		{"DRIZZ_" + secret + "=value"},
	}
	for _, environment := range cases {
		_, failure := configuration.New(environment).Load()
		if failure == nil {
			test.Fatalf("invalid configuration was accepted: %v", environment)
		}
		if strings.Contains(failure.Error(), secret) {
			test.Fatalf("a configuration error echoed a secret value: %v", failure)
		}
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	cases := []struct {
		name        string
		environment []string
		message     string
	}{
		{name: "unknown", environment: []string{"DRIZZ_UNKNOWN=value"}, message: "unknown Drizz setting"},
		{name: "misspelled", environment: []string{"DRIZZ_SENTRY_SAMPLE=0.5"}, message: "unknown Drizz setting"},
		{name: "notanumber", environment: []string{"DRIZZ_SENTRY_SAMPLE_RATE=NaN"}, message: "DRIZZ_SENTRY_SAMPLE_RATE"},
		{name: "positive", environment: []string{"DRIZZ_SENTRY_SAMPLE_RATE=+Inf"}, message: "DRIZZ_SENTRY_SAMPLE_RATE"},
		{name: "negative", environment: []string{"DRIZZ_SENTRY_SAMPLE_RATE=-Inf"}, message: "DRIZZ_SENTRY_SAMPLE_RATE"},
		{name: "level", environment: []string{"DRIZZ_LOG_LEVEL=trace"}, message: "DRIZZ_LOG_LEVEL must be one of"},
		{name: "endpoint", environment: []string{"DRIZZ_TELEMETRY_EXPORTER=OTLP"}, message: "is required"},
		{
			name:        "sample",
			environment: []string{"DRIZZ_SENTRY_SAMPLE_RATE=0"},
			message:     "DRIZZ_SENTRY_SAMPLE_RATE",
		},
		{
			name: "insecure",
			environment: []string{
				"DRIZZ_TELEMETRY_EXPORTER=OTLP",
				"DRIZZ_TELEMETRY_ENDPOINT=http://example.com",
			},
			message: "must use HTTPS",
		},
		{
			name: "credentials",
			environment: []string{
				"DRIZZ_TELEMETRY_EXPORTER=OTLP",
				"DRIZZ_TELEMETRY_ENDPOINT=https://user:secret@example.com",
			},
			message: "cannot contain credentials",
		},
	}
	for _, item := range cases {
		test.Run(item.name, func(test *testing.T) {
			test.Parallel()
			_, failure := configuration.New(item.environment).Load()
			if failure == nil || !strings.Contains(failure.Error(), item.message) {
				test.Fatalf("error = %v", failure)
			}
		})
	}
}
