package carrier_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/device/infrastructure/carrier"
)

type kit struct {
	test *testing.T
}

// probe builds a valid probe for some helper bytes: the bytes and their true digest.
func (kit kit) probe(bytes []byte) carrier.Probe {
	kit.test.Helper()
	sum := sha256.Sum256(bytes)
	return carrier.Probe{Bytes: bytes, Digest: hex.EncodeToString(sum[:])}
}

func (kit kit) read(path string) []byte {
	kit.test.Helper()
	raw, failure := os.ReadFile(path)
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return raw
}

func TestMaterializeReportsAbsentForPlaceholder(test *testing.T) {
	test.Setenv("XDG_CACHE_HOME", test.TempDir())
	test.Setenv("HOME", test.TempDir())

	made, failure := carrier.New()
	if failure != nil {
		test.Fatal(failure)
	}
	// The committed asset is the development placeholder, so a build that has not injected a real helper must report
	// no helper is carried rather than extract and run a stub.
	if _, failure := made.Materialize(); failure == nil {
		test.Fatal("a build carrying only the placeholder must report an absent helper")
	}
}

func TestMaterializeRunsRealHelper(test *testing.T) {
	test.Setenv("XDG_CACHE_HOME", test.TempDir())
	test.Setenv("HOME", test.TempDir())

	made, failure := carrier.New()
	if failure != nil {
		test.Fatal(failure)
	}
	path, failure := made.Materialize()
	if failure != nil {
		test.Skip("no real device helper is embedded in this build (placeholder)")
	}
	// A real helper was injected, so the extracted binary must run and answer the protocol.
	scope, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	command := exec.CommandContext(scope, path)
	command.Stdin = strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"health"}` + "\n")
	reply, failure := command.Output()
	if failure != nil {
		test.Fatalf("the extracted helper did not run: %v", failure)
	}
	if !strings.Contains(string(reply), `"protocol":1`) {
		test.Fatalf("helper reply = %q", reply)
	}
}

func TestDeliverWritesExecutable(test *testing.T) {
	test.Parallel()
	kit := kit{test: test}
	root := test.TempDir()

	path, failure := carrier.Wrap(root).Deliver(kit.probe([]byte("helper-bytes")))
	if failure != nil {
		test.Fatal(failure)
	}
	if string(kit.read(path)) != "helper-bytes" {
		test.Fatalf("extracted contents = %q", kit.read(path))
	}
	if runtime.GOOS != "windows" {
		info, failure := os.Stat(path)
		if failure != nil {
			test.Fatal(failure)
		}
		if info.Mode().Perm()&0o100 == 0 {
			test.Fatalf("the helper must be executable, mode = %v", info.Mode().Perm())
		}
	}
}

func TestDeliverReusesMatchingCopy(test *testing.T) {
	test.Parallel()
	kit := kit{test: test}
	root := test.TempDir()
	made := carrier.Wrap(root)

	first, failure := made.Deliver(kit.probe([]byte("same-bytes")))
	if failure != nil {
		test.Fatal(failure)
	}
	second, failure := made.Deliver(kit.probe([]byte("same-bytes")))
	if failure != nil {
		test.Fatal(failure)
	}
	if first != second {
		test.Fatalf("reuse must return the same path, got %q then %q", first, second)
	}
}

func TestDeliverReplacesTamperedCopy(test *testing.T) {
	test.Parallel()
	kit := kit{test: test}
	root := test.TempDir()
	made := carrier.Wrap(root)

	path, failure := made.Deliver(kit.probe([]byte("trusted-bytes")))
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := os.WriteFile(path, []byte("tampered"), 0o700); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := made.Deliver(kit.probe([]byte("trusted-bytes"))); failure != nil {
		test.Fatal(failure)
	}
	if string(kit.read(path)) != "trusted-bytes" {
		test.Fatal("a tampered helper must be re-extracted from the trusted bytes")
	}
}

func TestDeliverRefusesDigestMismatch(test *testing.T) {
	test.Parallel()
	root := test.TempDir()

	_, failure := carrier.Wrap(root).Deliver(carrier.Probe{Bytes: []byte("helper"), Digest: "deadbeef"})
	if failure == nil {
		test.Fatal("a helper whose bytes do not match the digest must be refused")
	}
	if entries, _ := os.ReadDir(root); len(entries) != 0 {
		test.Fatal("a refused helper must not be written")
	}
}

func TestDeliverRejectsSymlinkRoot(test *testing.T) {
	test.Parallel()
	kit := kit{test: test}
	base := test.TempDir()
	root := filepath.Join(base, "cache")
	if failure := os.Symlink(test.TempDir(), root); failure != nil {
		test.Fatal(failure)
	}

	_, failure := carrier.Wrap(root).Deliver(kit.probe([]byte("helper")))
	if failure == nil {
		test.Fatal("a symlinked cache root must be refused")
	}
}
