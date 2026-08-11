package sidecar_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/DrizzDev/platform/internal/device/infrastructure/sidecar"
)

func TestAbsent(test *testing.T) {
	test.Parallel()

	if _, ready := sidecar.New(nil).Locate(); ready {
		test.Fatal("an override was reported with none configured")
	}
}

func TestPinned(test *testing.T) {
	test.Parallel()

	handle, ready := sidecar.New([]string{
		"DRIZZ_DEVICE_SIDECAR=/opt/helper",
		"DRIZZ_DEVICE_DIGEST=abc123",
	}).Locate()
	if !ready {
		test.Fatal("a configured override was not reported")
	}
	if handle.Location != "/opt/helper" || handle.Digest != "abc123" {
		test.Fatalf("handle = %#v", handle)
	}
}

func TestComputed(test *testing.T) {
	test.Parallel()

	path := filepath.Join(test.TempDir(), "helper")
	if failure := os.WriteFile(path, []byte("helper-bytes"), 0o600); failure != nil {
		test.Fatal(failure)
	}
	handle, ready := sidecar.New([]string{"DRIZZ_DEVICE_SIDECAR=" + path}).Locate()
	if !ready {
		test.Fatal("a path-only override was not reported")
	}
	if handle.Location != path || len(handle.Digest) != 64 {
		test.Fatalf("handle = %#v", handle)
	}
}
