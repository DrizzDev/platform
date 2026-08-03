package sentry

func New(input Input) (Settings, error) {
	sample, failure := input.parse()
	if failure != nil {
		return Settings{}, failure
	}
	settings := Settings{
		dsn:         input.DSN,
		environment: input.Environment,
		sample:      sample,
	}
	if failure := settings.validate(); failure != nil {
		return Settings{}, failure
	}
	return settings, nil
}
