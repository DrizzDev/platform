package credential

import (
	"errors"
	"slices"
	"time"

	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

const (
	limit  = 512
	secret = 512
	mark   = 43
)

type Record struct {
	issuer string
	client string
	handle string

	subject subject.Subject
	session session.Session

	refresh []byte
	method  method.Method

	issued time.Time
	expiry time.Time

	revision uint64
	schema   uint32
}

type Input struct {
	Issuer string
	Client string
	Handle string

	Subject subject.Subject
	Session session.Session

	Refresh []byte
	Method  method.Method

	Issued time.Time
	Expiry time.Time

	Revision uint64
	Schema   uint32
}

func New(input Input) (Record, error) {
	record := Record{
		issuer:   input.Issuer,
		client:   input.Client,
		handle:   input.Handle,
		subject:  input.Subject,
		session:  input.Session,
		method:   input.Method,
		issued:   input.Issued,
		expiry:   input.Expiry,
		schema:   input.Schema,
		revision: input.Revision,
		refresh:  slices.Clone(input.Refresh),
	}
	if failure := record.validate(); failure != nil {
		return Record{}, failure
	}
	return record, nil
}

func (record Record) Issuer() string {
	return record.issuer
}

func (record Record) Client() string {
	return record.client
}

func (record Record) Handle() string {
	return record.handle
}

func (record Record) Subject() subject.Subject {
	return record.subject
}

func (record Record) Session() session.Session {
	return record.session
}

func (record Record) Method() method.Method {
	return record.method
}

func (record Record) Refresh() []byte {
	return slices.Clone(record.refresh)
}

func (record Record) Issued() time.Time {
	return record.issued
}

func (record Record) Expiry() time.Time {
	return record.expiry
}

func (record Record) Revision() uint64 {
	return record.revision
}

func (record Record) Schema() uint32 {
	return record.schema
}

func (record Record) validate() error {
	if failure := record.identity(); failure != nil {
		return failure
	}
	if failure := record.principal(); failure != nil {
		return failure
	}
	return record.lifecycle()
}

func (record Record) identity() error {
	switch {
	case record.issuer == "":
		return errors.New("credential issuer is required")
	case len(record.issuer) > limit:
		return errors.New("credential issuer exceeds the maximum length")
	case record.client == "":
		return errors.New("credential client is required")
	case len(record.client) > limit:
		return errors.New("credential client exceeds the maximum length")
	case record.handle == "":
		return errors.New("credential handle is required")
	case len(record.handle) > mark:
		return errors.New("credential handle exceeds the maximum length")
	default:
		return nil
	}
}

func (record Record) principal() error {
	switch {
	case record.subject.String() == "":
		return errors.New("credential subject is required")
	case record.session.String() == "":
		return errors.New("credential session is required")
	case !record.method.Valid():
		return errors.New("credential method is invalid")
	case len(record.refresh) == 0:
		return errors.New("credential refresh is required")
	case len(record.refresh) > secret:
		return errors.New("credential refresh exceeds the maximum length")
	default:
		return nil
	}
}

func (record Record) lifecycle() error {
	switch {
	case record.issued.IsZero():
		return errors.New("credential issue time is required")
	case record.expiry.IsZero():
		return errors.New("credential expiry is required")
	case !record.expiry.After(record.issued):
		return errors.New("credential expiry must be after the issue time")
	case record.revision < 1:
		return errors.New("credential revision must be positive")
	case record.schema < 1:
		return errors.New("credential schema must be positive")
	default:
		return nil
	}
}
