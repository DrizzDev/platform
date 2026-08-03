package scope_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/DrizzDev/platform/internal/platform/logging/scope"
)

func TestBind(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&output, nil))
	logger = scope.New(scope.Input{
		Request:     "request_123",
		Correlation: "correlation_123",
		Session:     "session_123",
	}).Bind(logger)
	logger.Info("ready")

	var record map[string]any
	if failure := json.Unmarshal(output.Bytes(), &record); failure != nil {
		test.Fatal(failure)
	}
	if record["request_id"] != "request_123" ||
		record["correlation_id"] != "correlation_123" ||
		record["session_id"] != "session_123" {
		test.Fatalf("record = %#v", record)
	}
}

func TestOmit(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&output, nil))
	scope.New(scope.Input{}).Bind(logger).Info("ready")

	var record map[string]any
	if failure := json.Unmarshal(output.Bytes(), &record); failure != nil {
		test.Fatal(failure)
	}
	for _, key := range []string{"request_id", "correlation_id", "session_id"} {
		if _, found := record[key]; found {
			test.Fatalf("%s must be omitted", key)
		}
	}
}
