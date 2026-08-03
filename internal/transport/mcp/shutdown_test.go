package mcp_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/transport/mcp"
)

func TestShutdown(test *testing.T) {
	test.Parallel()

	reader, _ := io.Pipe()
	if !mcp.Terminate(reader) {
		test.Fatal("closing the connection did not terminate the blocked reader")
	}
}

func TestCleanup(test *testing.T) {
	test.Parallel()

	first, second := mcp.Discard(unclosable{Reader: strings.NewReader("")})
	if first == nil || second == nil {
		test.Fatalf("a failed input close was not reported: %v then %v", first, second)
	}
	if first.Error() != second.Error() {
		test.Fatalf("repeated close was not idempotent: %v then %v", first, second)
	}
}

type unclosable struct {
	io.Reader
}

func (unclosable) Close() error {
	return errors.New("input close failed")
}
