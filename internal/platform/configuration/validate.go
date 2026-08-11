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
				"DRIZZ_SENTRY_DSN, DRIZZ_SENTRY_SAMPLE_RATE, DRIZZ_SENTRY_ENVIRONMENT, " +
				"DRIZZ_AUTH0_ISSUER, DRIZZ_AUTH0_CLIENT, DRIZZ_AUTH0_AUDIENCE, " +
				"DRIZZ_AUTH0_REDIRECT, DRIZZ_AUTH0_SCOPES, DRIZZ_SESSION, DRIZZ_CLOUD")
		}
	}
	return nil
}

func (loader Loader) reject(key string) bool {
	if !strings.HasPrefix(key, prefix) {
		return false
	}
	if strings.HasPrefix(key, reserved) {
		return false
	}
	switch key {
	case level, exporter, endpoint, dsn, sample, stage,
		issuer, client, audience, redirect, scopes, session, cloud:
		return false
	default:
		return true
	}
}
