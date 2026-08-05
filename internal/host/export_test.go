package host

import (
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/platform/configuration/identity"
)

// Authority exposes the post-authentication gate builder for tests, keyed by the
// configured cloud base URL.
func Authority(cloud string) (login.Authority, error) {
	return foundation{}.authority(kit{settings: identity.New(identity.Input{Cloud: cloud})})
}
