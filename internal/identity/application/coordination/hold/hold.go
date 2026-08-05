package hold

import (
	"errors"
	"time"

	"github.com/DrizzDev/platform/internal/identity/domain/session"
)

type Hold struct {
	owner   string
	moment  time.Time
	window  time.Duration
	session session.Session
}

type Input struct {
	Owner   string
	Moment  time.Time
	Window  time.Duration
	Session session.Session
}

func New(input Input) (Hold, error) {
	hold := Hold{session: input.Session, owner: input.Owner, moment: input.Moment, window: input.Window}
	if failure := hold.validate(); failure != nil {
		return Hold{}, failure
	}
	return hold, nil
}

func (hold Hold) Session() session.Session {
	return hold.session
}

func (hold Hold) Owner() string {
	return hold.owner
}

func (hold Hold) Moment() time.Time {
	return hold.moment
}

func (hold Hold) Window() time.Duration {
	return hold.window
}

func (hold Hold) Deadline() time.Time {
	return hold.moment.Add(hold.window)
}

func (hold Hold) validate() error {
	switch {
	case hold.session.String() == "":
		return errors.New("hold session is required")
	case hold.owner == "":
		return errors.New("hold owner is required")
	case hold.moment.IsZero():
		return errors.New("hold moment is required")
	case hold.window <= 0:
		return errors.New("hold window must be positive")
	default:
		return nil
	}
}
