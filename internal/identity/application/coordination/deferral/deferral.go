package deferral

import (
	"errors"
	"time"
)

type Deferral struct {
	key  string
	next time.Time
}

type Input struct {
	Key  string
	Next time.Time
}

func New(input Input) (Deferral, error) {
	deferral := Deferral{key: input.Key, next: input.Next}
	if failure := deferral.validate(); failure != nil {
		return Deferral{}, failure
	}
	return deferral, nil
}

func (deferral Deferral) Key() string {
	return deferral.key
}

func (deferral Deferral) Next() time.Time {
	return deferral.next
}

func (deferral Deferral) validate() error {
	switch {
	case deferral.key == "":
		return errors.New("deferral key is required")
	case deferral.next.IsZero():
		return errors.New("deferral next time is required")
	default:
		return nil
	}
}
