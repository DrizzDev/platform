package bridge

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"

	"github.com/DrizzDev/platform/internal/device/application/control"
)

// read runs a read-only request with bounded retry: only a transient transport outage is retried, with backoff and
// jitter, never a device error (REL-013).
func (driver *Driver) read(scope context.Context, request Request) (Response, error) {
	var failure error

	for attempt := range attempts {
		response, cause := driver.channel.Invoke(scope, request)
		if cause == nil {
			return response, nil
		}
		failure = cause
		if !errors.Is(cause, Unavailable{}) || attempt == attempts-1 {
			break
		}
		if driver.pause(scope, attempt) != nil {
			break
		}
	}
	return Response{}, failure
}

func (driver *Driver) pause(scope context.Context, attempt int) error {
	span := lull << attempt
	timer := time.NewTimer(span + time.Duration(rand.Int64N(int64(span))))

	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-scope.Done():
		return scope.Err()
	}
}

// fault maps a transport failure to the device port's typed error; a context error passes through so the flow
// renders cancellation.
func (driver *Driver) fault(cause error) error {
	var problem Fault

	switch {
	case errors.As(cause, &problem):
		if mapped, found := kinds[problem.Kind()]; found {
			return mapped
		}
		return control.Unavailable{}
	case errors.Is(cause, Mismatch{}):
		return control.Incompatible{}
	case errors.Is(cause, Oversize{}):
		return control.Failed{}
	case errors.Is(cause, Unavailable{}):
		return control.Unavailable{}
	default:
		return cause
	}
}
