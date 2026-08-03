package telemetry

import (
	"context"
	"errors"

	"github.com/DrizzDev/platform/internal/platform/configuration/telemetry"
)

func New(scope context.Context, options Options) (Provider, error) {
	if failure := options.validate(); failure != nil {
		return Provider{}, failure
	}
	switch options.Settings.Exporter() {
	case telemetry.None:
		return options.noop(), nil
	case telemetry.OTLP:
		return options.otlp(scope)
	default:
		return Provider{}, errors.New("telemetry exporter is unsupported")
	}
}
