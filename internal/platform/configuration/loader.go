package configuration

import (
	"strings"

	"github.com/DrizzDev/platform/internal/platform/configuration/logging"
	"github.com/DrizzDev/platform/internal/platform/configuration/reporting"
	"github.com/DrizzDev/platform/internal/platform/configuration/reporting/sentry"
	"github.com/DrizzDev/platform/internal/platform/configuration/telemetry"
)

type Loader struct {
	values map[string]string
}

func New(environment []string) Loader {
	values := make(map[string]string, len(environment))
	for _, entry := range environment {
		key, value, found := strings.Cut(entry, "=")
		if found {
			values[key] = value
		}
	}
	return Loader{values: values}
}

func (loader Loader) Load() (Settings, error) {
	if failure := loader.validate(); failure != nil {
		return Settings{}, failure
	}
	logging, failure := logging.New(loader.values[level])
	if failure != nil {
		return Settings{}, failure
	}
	reporting, failure := reporting.New(reporting.Input{
		Sentry: sentry.Input{
			DSN:         loader.values[dsn],
			Environment: loader.values[stage],
			Sample:      loader.values[sample],
		},
	})
	if failure != nil {
		return Settings{}, failure
	}
	telemetry, failure := telemetry.New(telemetry.Input{
		Exporter: loader.values[exporter],
		Endpoint: loader.values[endpoint],
	})
	if failure != nil {
		return Settings{}, failure
	}
	return Settings{logging: logging, reporting: reporting, telemetry: telemetry}, nil
}
