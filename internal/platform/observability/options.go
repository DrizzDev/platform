package observability

import (
	"io"

	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/platform/configuration"
)

type Options struct {
	Output   io.Writer
	Build    build.Info
	Settings configuration.Settings
}
