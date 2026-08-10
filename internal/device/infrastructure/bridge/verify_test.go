package bridge_test

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/DrizzDev/platform/internal/device/infrastructure/bridge"
)

type vault struct {
	test *testing.T
}

func (vault vault) binary(content []byte) (string, string) {
	vault.test.Helper()
	path := filepath.Join(vault.test.TempDir(), "sidecar")
	if failure := os.WriteFile(path, content, 0o600); failure != nil {
		vault.test.Fatal(failure)
	}
	sum := sha256.Sum256(content)
	return path, hex.EncodeToString(sum[:])
}

func TestVerify(test *testing.T) {
	test.Parallel()

	path, digest := vault{test: test}.binary([]byte("sidecar-bytes"))
	if failure := bridge.Verify(bridge.Options{Location: path, Digest: digest}); failure != nil {
		test.Fatalf("a valid binary was rejected: %v", failure)
	}
}

func TestVerifyRejects(test *testing.T) {
	test.Parallel()

	kit := vault{test: test}
	path, digest := kit.binary([]byte("real"))

	if bridge.Verify(bridge.Options{Location: path, Digest: "deadbeef"}) == nil {
		test.Fatal("a mismatched digest was accepted")
	}
	if bridge.Verify(bridge.Options{Location: filepath.Join(test.TempDir(), "absent"), Digest: digest}) == nil {
		test.Fatal("a missing binary was accepted")
	}

	writable, sum := kit.binary([]byte("open"))
	if failure := os.Chmod(writable, 0o666); failure != nil {
		test.Fatal(failure)
	}
	if bridge.Verify(bridge.Options{Location: writable, Digest: sum}) == nil {
		test.Fatal("a world-writable binary was accepted")
	}
}
