package publication

import (
	"errors"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
)

type Publication struct {
	key      string
	revision uint64
	moment   time.Time
	expected epoch.Epoch
	session  session.Session
}

type Input struct {
	Key      string
	Revision uint64
	Moment   time.Time
	Expected epoch.Epoch
	Session  session.Session
}

func New(input Input) (Publication, error) {
	publication := Publication{
		key:      input.Key,
		moment:   input.Moment,
		session:  input.Session,
		revision: input.Revision,
		expected: input.Expected,
	}
	if failure := publication.validate(); failure != nil {
		return Publication{}, failure
	}
	return publication, nil
}

func (publication Publication) Session() session.Session {
	return publication.session
}

func (publication Publication) Expected() epoch.Epoch {
	return publication.expected
}

func (publication Publication) Key() string {
	return publication.key
}

func (publication Publication) Revision() uint64 {
	return publication.revision
}

func (publication Publication) Moment() time.Time {
	return publication.moment
}

func (publication Publication) validate() error {
	switch {
	case publication.session.String() == "":
		return errors.New("publication session is required")
	case publication.key == "":
		return errors.New("publication key is required")
	case publication.revision < 1:
		return errors.New("publication revision must be positive")
	case publication.moment.IsZero():
		return errors.New("publication moment is required")
	default:
		return nil
	}
}
