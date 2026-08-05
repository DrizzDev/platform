package interactive

import "github.com/DrizzDev/platform/internal/identity/application/login"

type Options struct {
	Authorization login.Authorization
	Browser       login.Browser
	Random        login.Random
}
