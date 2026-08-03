package sentry

import (
	"errors"
	"net/url"
)

const environment = 63

func (settings Settings) validate() error {
	if failure := settings.stage(); failure != nil {
		return failure
	}
	if !settings.Enabled() {
		return nil
	}
	address, failure := url.Parse(settings.dsn)
	if failure != nil || address.Scheme != "https" || address.Host == "" {
		return errors.New("DRIZZ_SENTRY_DSN must be an absolute HTTPS URL")
	}
	return nil
}

func (settings Settings) stage() error {
	if settings.environment == "" {
		return nil
	}
	if len(settings.environment) > environment {
		return errors.New("DRIZZ_SENTRY_ENVIRONMENT must be at most 63 characters")
	}
	for index := range len(settings.environment) {
		if !settings.permitted(settings.environment[index]) {
			return errors.New("DRIZZ_SENTRY_ENVIRONMENT must use letters, digits, '.', '_', or '-'")
		}
	}
	return nil
}

func (settings Settings) permitted(symbol byte) bool {
	switch {
	case symbol >= 'a' && symbol <= 'z', symbol >= 'A' && symbol <= 'Z':
		return true
	case symbol >= '0' && symbol <= '9':
		return true
	default:
		return symbol == '.' || symbol == '_' || symbol == '-'
	}
}
