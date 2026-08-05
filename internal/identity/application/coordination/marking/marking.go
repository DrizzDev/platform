package marking

import (
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/attempt"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
)

type Marking struct {
	session session.Session
	attempt attempt.Attempt
}

type Input struct {
	Session session.Session
	Attempt attempt.Attempt
}

func New(input Input) (Marking, error) {
	marking := Marking{session: input.Session, attempt: input.Attempt}
	if failure := marking.validate(); failure != nil {
		return Marking{}, failure
	}
	return marking, nil
}

func (marking Marking) Session() session.Session {
	return marking.session
}

func (marking Marking) Attempt() attempt.Attempt {
	return marking.attempt
}

func (marking Marking) validate() error {
	if marking.session.String() == "" {
		return errors.New("marking session is required")
	}
	if marking.attempt.Revision() < 1 {
		return errors.New("marking attempt is required")
	}
	return nil
}
