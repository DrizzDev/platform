package reporting

import "github.com/DrizzDev/platform/internal/platform/configuration/reporting/sentry"

type Settings struct {
	sentry sentry.Settings
}

func (settings Settings) Sentry() sentry.Settings {
	return settings.sentry
}
