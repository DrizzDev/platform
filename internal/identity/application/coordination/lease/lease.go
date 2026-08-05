package lease

import (
	"errors"
	"time"
)

type Lease struct {
	owner  string
	expiry time.Time
}

type Input struct {
	Owner  string
	Expiry time.Time
}

func New(input Input) (Lease, error) {
	lease := Lease{owner: input.Owner, expiry: input.Expiry}
	if failure := lease.validate(); failure != nil {
		return Lease{}, failure
	}
	return lease, nil
}

func (lease Lease) Owner() string {
	return lease.owner
}

func (lease Lease) Expiry() time.Time {
	return lease.expiry
}

func (lease Lease) Held(now time.Time) bool {
	return now.Before(lease.expiry)
}

func (lease Lease) validate() error {
	switch {
	case lease.owner == "":
		return errors.New("lease owner is required")
	case lease.expiry.IsZero():
		return errors.New("lease expiry is required")
	default:
		return nil
	}
}
