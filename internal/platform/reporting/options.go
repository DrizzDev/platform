package reporting

import (
	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/platform/configuration/reporting"
)

type Options struct {
	Build    build.Info
	Settings reporting.Settings
}
