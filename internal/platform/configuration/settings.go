package configuration

import (
	"github.com/DrizzDev/platform/internal/platform/configuration/identity"
	"github.com/DrizzDev/platform/internal/platform/configuration/logging"
	"github.com/DrizzDev/platform/internal/platform/configuration/reporting"
	"github.com/DrizzDev/platform/internal/platform/configuration/telemetry"
)

type Settings struct {
	logging   logging.Settings
	reporting reporting.Settings
	telemetry telemetry.Settings
	identity  identity.Settings
}

func (settings Settings) Identity() identity.Settings {
	return settings.identity
}

func (settings Settings) Logging() logging.Settings {
	return settings.logging
}

func (settings Settings) Telemetry() telemetry.Settings {
	return settings.telemetry
}

func (settings Settings) Reporting() reporting.Settings {
	return settings.reporting
}
