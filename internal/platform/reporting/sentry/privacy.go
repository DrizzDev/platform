package sentry

import sdk "github.com/getsentry/sentry-go"

func (Options) privacy() *sdk.DataCollection {
	disabled := &sdk.KeyValueCollectionBehavior{Mode: sdk.CollectionOff}
	return &sdk.DataCollection{
		Cookies:  disabled,
		UserInfo: sdk.Set(false),
		HTTPHeaders: &sdk.HeaderCollectionConfig{
			Request:  disabled,
			Response: disabled,
		},
		QueryParams: disabled,
		HTTPBodies:  []sdk.BodyType{},
	}
}
