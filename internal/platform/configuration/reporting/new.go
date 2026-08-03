package reporting

import "github.com/DrizzDev/platform/internal/platform/configuration/reporting/sentry"

func New(input Input) (Settings, error) {
	settings, failure := sentry.New(input.Sentry)
	if failure != nil {
		return Settings{}, failure
	}
	return Settings{sentry: settings}, nil
}
