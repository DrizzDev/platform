package login

import "context"

// Browser is the login-owned port over the system browser and loopback
// listener. Prompt opens the redirect and returns the single captured callback.
type Browser interface {
	Prompt(context.Context, Redirect) (Callback, error)
}

type Callback struct {
	Code  string
	State string
}
