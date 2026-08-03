package logging

import (
	"io"

	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/platform/configuration/logging"
)

type Options struct {
	Output   io.Writer
	Build    build.Info
	Settings logging.Settings
}
