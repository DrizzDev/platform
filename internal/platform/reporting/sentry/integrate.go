package sentry

import (
	sdk "github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
)

func (Options) integrate(integrations []sdk.Integration) []sdk.Integration {
	excluded := map[string]bool{"Environment": true, "Modules": true}
	kept := make([]sdk.Integration, 0, len(integrations)+1)
	for _, integration := range integrations {
		if !excluded[integration.Name()] {
			kept = append(kept, integration)
		}
	}
	return append(kept, sentryotel.NewOtelIntegration())
}
