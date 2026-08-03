package sentry

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	sdk "github.com/getsentry/sentry-go"
)

type Provider struct {
	handler slog.Handler
	client  *sdk.Client
	once    *sync.Once
	outcome *error
}

func (provider Provider) Handler() slog.Handler {
	return provider.handler
}

func (provider Provider) Close(scope context.Context) error {
	if provider.client == nil {
		return nil
	}
	provider.once.Do(func() {
		sent := provider.client.FlushWithContext(scope)
		provider.client.Close()
		if !sent {
			*provider.outcome = errors.New("sentry flush exceeded the shutdown deadline")
		}
	})
	return *provider.outcome
}
