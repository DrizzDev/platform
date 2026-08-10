package policy

import (
	"errors"
	"time"
)

// Policy is the handling rule for one data category: a byte budget, a retention window, whether it may be uploaded,
// and whether it must be redacted before storage.
type Policy struct {
	limit     int
	retention time.Duration
	upload    bool
	redaction bool
}

type Input struct {
	Limit     int
	Retention time.Duration
	Upload    bool
	Redaction bool
}

func New(input Input) (Policy, error) {
	policy := Policy{limit: input.Limit, retention: input.Retention, upload: input.Upload, redaction: input.Redaction}
	if failure := policy.validate(); failure != nil {
		return Policy{}, failure
	}
	return policy, nil
}

func (policy Policy) Limit() int {
	return policy.limit
}

func (policy Policy) Retention() time.Duration {
	return policy.retention
}

func (policy Policy) Upload() bool {
	return policy.upload
}

func (policy Policy) Redaction() bool {
	return policy.redaction
}

func (policy Policy) validate() error {
	switch {
	case policy.limit <= 0:
		return errors.New("policy limit must be positive")
	case policy.retention <= 0:
		return errors.New("policy retention must be positive")
	default:
		return nil
	}
}
