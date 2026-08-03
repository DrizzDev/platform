package telemetry_test

import (
	"context"
	"testing"

	"github.com/DrizzDev/platform/internal/platform/build"
	configuration "github.com/DrizzDev/platform/internal/platform/configuration/telemetry"
	"github.com/DrizzDev/platform/internal/platform/telemetry"
)

func TestDisabled(test *testing.T) {
	test.Parallel()

	settings, failure := configuration.New(configuration.Input{})
	if failure != nil {
		test.Fatal(failure)
	}
	provider, failure := telemetry.New(
		context.Background(),
		telemetry.Options{Settings: settings, Build: build.Read()},
	)
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := provider.Close(context.Background()); failure != nil {
		test.Fatal(failure)
	}
}

func TestIdentity(test *testing.T) {
	test.Parallel()

	_, failure := telemetry.New(context.Background(), telemetry.Options{})
	if failure == nil {
		test.Fatal("missing identity was accepted")
	}
}
