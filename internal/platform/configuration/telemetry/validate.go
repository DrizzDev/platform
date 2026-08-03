package telemetry

import (
	"errors"
	"net"
	"net/url"
)

func (settings Settings) validate() error {
	switch settings.exporter {
	case None:
		if settings.endpoint != "" {
			return errors.New("DRIZZ_TELEMETRY_ENDPOINT requires the OTLP exporter")
		}
		return nil
	case OTLP:
		if settings.endpoint == "" {
			return errors.New("DRIZZ_TELEMETRY_ENDPOINT is required for the OTLP exporter")
		}
		return settings.address()
	default:
		return errors.New("DRIZZ_TELEMETRY_EXPORTER must be none or otlp")
	}
}

func (settings Settings) address() error {
	endpoint, failure := url.Parse(settings.endpoint)
	if failure != nil || endpoint.Host == "" {
		return errors.New("DRIZZ_TELEMETRY_ENDPOINT must be an absolute URL")
	}
	if endpoint.User != nil || endpoint.RawQuery != "" || endpoint.Fragment != "" {
		return errors.New("DRIZZ_TELEMETRY_ENDPOINT cannot contain credentials, queries, or fragments")
	}
	if endpoint.Scheme == "https" {
		return nil
	}
	host := endpoint.Hostname()
	if endpoint.Scheme == "http" && (host == "localhost" || net.ParseIP(host).IsLoopback()) {
		return nil
	}
	return errors.New("DRIZZ_TELEMETRY_ENDPOINT must use HTTPS unless it targets this computer")
}
