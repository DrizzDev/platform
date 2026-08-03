package configuration

import (
	"errors"
	"strings"
)

func (loader Loader) validate() error {
	for key := range loader.values {
		if loader.reject(key) {
			return errors.New("unknown Drizz setting; supported settings are " +
				"DRIZZ_LOG_LEVEL, DRIZZ_TELEMETRY_EXPORTER, DRIZZ_TELEMETRY_ENDPOINT, " +
				"DRIZZ_SENTRY_DSN, DRIZZ_SENTRY_SAMPLE_RATE, DRIZZ_SENTRY_ENVIRONMENT")
		}
	}
	return nil
}

func (loader Loader) reject(key string) bool {
	if !strings.HasPrefix(key, prefix) {
		return false
	}
	switch key {
	case level, exporter, endpoint, dsn, sample, stage:
		return false
	default:
		return true
	}
}
