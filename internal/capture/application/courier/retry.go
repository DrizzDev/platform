package courier

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"
)

// retry runs a remote upload with bounded attempts. A terminal rejection or a cancelled context stops immediately; a
// transient failure backs off with jitter and tries again.
func (courier Courier) retry(scope context.Context, work func(context.Context) error) error {
	var failure error

	for attempt := range attempts {
		failure = work(scope)
		if failure == nil {
			return nil
		}
		var stop terminal
		if errors.As(failure, &stop) || errors.Is(failure, context.Canceled) || attempt == attempts-1 {
			return failure
		}
		if courier.pause(scope, attempt) != nil {
			return failure
		}
	}
	return failure
}

func (Courier) pause(scope context.Context, attempt int) error {
	span := min(lull<<attempt, ceiling)
	timer := time.NewTimer(span + time.Duration(rand.Int64N(int64(span))))

	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-scope.Done():
		return scope.Err()
	}
}
