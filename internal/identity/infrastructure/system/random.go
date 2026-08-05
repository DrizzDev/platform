package system

import (
	"crypto/rand"

	"github.com/DrizzDev/platform/internal/identity/application/login"
)

var _ login.Random = Random{}

// Random draws cryptographic randomness from the operating system (SEC-005).
type Random struct{}

func (Random) Bytes(size int) ([]byte, error) {
	buffer := make([]byte, size)
	if _, failure := rand.Read(buffer); failure != nil {
		return nil, failure
	}
	return buffer, nil
}
