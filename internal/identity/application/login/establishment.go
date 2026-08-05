package login

import "context"

// Establishment is the login-owned seam over the first stage of a sign-in: it
// proves who the caller is and yields a validated provider token. The browser
// and device front-ends are interchangeable implementations, and the flow
// publishes the result identically for either.
type Establishment interface {
	Establish(context.Context) (Token, error)
}
