package session

import (
	"context"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
)

// Refresh exchanges the stored refresh token at the provider for a rotated
// credential. It runs once per renewal and is never retried: a rotating refresh
// token may already have been consumed, so any failure is treated as uncertain.
type Refresh interface {
	Renew(context.Context, credential.Record) (Renewal, error)
}

// Renewal is the validated provider result of a refresh: the rotated refresh
// token, the access token for cloud calls, and the access-token expiry that
// dates the next renewal.
type Renewal struct {
	Refresh []byte
	Access  []byte
	Expiry  time.Time
}
