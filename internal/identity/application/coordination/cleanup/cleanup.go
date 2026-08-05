package cleanup

import (
	"errors"
	"time"
)

type Record struct {
	key      string
	reason   Reason
	state    State
	attempts uint
	next     time.Time
	deadline time.Time
	created  time.Time
}

type Input struct {
	Key      string
	Reason   Reason
	State    State
	Attempts uint
	Next     time.Time
	Deadline time.Time
	Created  time.Time
}

func New(input Input) (Record, error) {
	record := Record{
		key:      input.Key,
		next:     input.Next,
		state:    input.State,
		reason:   input.Reason,
		created:  input.Created,
		attempts: input.Attempts,
		deadline: input.Deadline,
	}
	if failure := record.validate(); failure != nil {
		return Record{}, failure
	}
	return record, nil
}

func (record Record) Key() string {
	return record.key
}

func (record Record) Reason() Reason {
	return record.reason
}

func (record Record) State() State {
	return record.state
}

func (record Record) Attempts() uint {
	return record.attempts
}

func (record Record) Next() time.Time {
	return record.next
}

func (record Record) Deadline() time.Time {
	return record.deadline
}

func (record Record) Created() time.Time {
	return record.created
}

func (record Record) Blocked() bool {
	return record.state == Blocked
}

func (record Record) validate() error {
	switch {
	case record.key == "":
		return errors.New("cleanup key is required")

	case !record.reason.Valid():
		return errors.New("cleanup reason is invalid")

	case !record.state.Valid():
		return errors.New("cleanup state is invalid")

	case record.created.IsZero():
		return errors.New("cleanup creation time is required")

	case record.deadline.IsZero():
		return errors.New("cleanup deadline is required")

	default:
		return nil
	}
}
