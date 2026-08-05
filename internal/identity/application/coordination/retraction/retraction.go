package retraction

import (
	"errors"
	"time"

	"github.com/DrizzDev/platform/internal/identity/domain/session"
)

type Retraction struct {
	session session.Session
	moment  time.Time
}

type Input struct {
	Session session.Session
	Moment  time.Time
}

func New(input Input) (Retraction, error) {
	retraction := Retraction{session: input.Session, moment: input.Moment}
	if failure := retraction.validate(); failure != nil {
		return Retraction{}, failure
	}
	return retraction, nil
}

func (retraction Retraction) Session() session.Session {
	return retraction.session
}

func (retraction Retraction) Moment() time.Time {
	return retraction.moment
}

func (retraction Retraction) validate() error {
	switch {
	case retraction.session.String() == "":
		return errors.New("retraction session is required")
	case retraction.moment.IsZero():
		return errors.New("retraction moment is required")
	default:
		return nil
	}
}
