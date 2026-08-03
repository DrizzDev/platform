package logging

import (
	"log/slog"
	"strings"
)

type redactor struct{}

func (redactor redactor) replace(input replacement) slog.Attr {
	if redactor.sensitive(input.attribute.Key) {
		return slog.String(input.attribute.Key, redacted)
	}
	for _, group := range input.groups {
		if redactor.sensitive(group) {
			return slog.String(input.attribute.Key, redacted)
		}
	}
	return input.attribute
}

func (redactor redactor) sensitive(name string) bool {
	switch strings.ToLower(name) {
	case authorization,
		cookie,
		credential,
		credentials,
		password,
		secret,
		client,
		token,
		access,
		refresh,
		idtoken,
		private,
		api,
		apikey:
		return true
	default:
		return false
	}
}
