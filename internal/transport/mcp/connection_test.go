package mcp_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/transport/mcp"
)

type expectation struct {
	name     string
	input    string
	rejected bool
	clean    bool
}

func TestFraming(test *testing.T) {
	test.Parallel()

	request := `{"jsonrpc":"2.0","id":1,"method":"ping"}`
	cases := []expectation{
		{name: "valid", input: request + "\n"},
		{name: "blank", input: "\n\n" + request + "\n"},
		{name: "oversize", input: strings.Repeat("a", (1<<20)+16) + "\n", rejected: true},
		{name: "multiline", input: "{\n\"jsonrpc\":\"2.0\"\n}\n", rejected: true},
		{name: "truncated", input: `{"jsonrpc":"2.0","id":1,"me`, rejected: true},
		{name: "trailing", input: request + " trailing\n", rejected: true},
		{name: "double", input: request + request + "\n", rejected: true},
		{name: "garbage", input: "not json at all\n", rejected: true},
		{name: "empty", input: "", clean: true},
	}
	for _, item := range cases {
		test.Run(item.name, func(test *testing.T) {
			test.Parallel()
			item.verify(test)
		})
	}
}

func (expectation expectation) verify(test *testing.T) {
	test.Helper()
	message, failure := mcp.Read(strings.NewReader(expectation.input))
	switch {
	case expectation.rejected:
		if !mcp.Rejected(failure) {
			test.Fatalf("input %q was not rejected: %v", expectation.input, failure)
		}
	case expectation.clean:
		if !errors.Is(failure, io.EOF) {
			test.Fatalf("input %q was not a clean end: %v", expectation.input, failure)
		}
	default:
		if failure != nil || message == nil {
			test.Fatalf("valid input failed: message=%v failure=%v", message, failure)
		}
	}
}
