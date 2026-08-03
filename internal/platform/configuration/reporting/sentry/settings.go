package sentry

type Settings struct {
	dsn         string
	environment string
	sample      float64
}

func (settings Settings) Enabled() bool {
	return settings.dsn != ""
}

func (settings Settings) DSN() string {
	return settings.dsn
}

func (settings Settings) Environment() string {
	return settings.environment
}

func (settings Settings) Sample() float64 {
	return settings.sample
}
