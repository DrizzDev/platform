package auth0

import (
	"net/http"

	"github.com/DrizzDev/platform/internal/identity/domain/method"
)

type Options struct {
	Agent    *http.Client
	Issuer   string
	Client   string
	Audience string
	Redirect string
	Method   method.Method
	Scopes   []string
}
