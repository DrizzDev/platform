package sidecar

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
)

const (
	pin    = "DRIZZ_DEVICE_DIGEST"
	binary = "DRIZZ_DEVICE_SIDECAR"
)

// Handle locates the compiled device helper and the digest its bytes must match before it is run.
type Handle struct {
	Location string
	Digest   string
}

// Resolver reads the device helper location from an environment. The override it reads is for tests and continuous
// integration only; production resolves the helper from the copy carried inside the application.
type Resolver struct {
	environment []string
}

func New(environment []string) Resolver {
	return Resolver{environment: environment}
}

// Locate reports the helper handle and whether an override is configured. With none configured the device capability
// is unprepared. A digest supplied beside the path pins the bytes; when it is omitted the digest is computed from the
// file for a local run.
func (resolver Resolver) Locate() (Handle, bool) {
	location := resolver.value(binary)
	if location == "" {
		return Handle{}, false
	}
	digest := resolver.value(pin)
	if digest == "" {
		computed, failure := resolver.fingerprint(location)
		if failure != nil {
			return Handle{}, false
		}
		digest = computed
	}
	return Handle{Location: location, Digest: digest}, true
}

func (resolver Resolver) value(key string) string {
	for _, entry := range resolver.environment {
		if value, found := strings.CutPrefix(entry, key+"="); found {
			return value
		}
	}
	return ""
}

func (resolver Resolver) fingerprint(path string) (string, error) {
	content, failure := os.ReadFile(path)
	if failure != nil {
		return "", failure
	}
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:]), nil
}
