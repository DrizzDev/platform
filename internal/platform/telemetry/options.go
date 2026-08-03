package telemetry

import (
	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/platform/configuration/telemetry"
)

type Options struct {
	Build    build.Info
	Settings telemetry.Settings
}
