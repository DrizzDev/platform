package logging_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/platform/configuration/logging"
	platform "github.com/DrizzDev/platform/internal/platform/logging"
)

func TestRecord(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	settings, failure := logging.New("debug")
	if failure != nil {
		test.Fatal(failure)
	}
	logger, failure := platform.New(platform.Options{
		Output:   &output,
		Settings: settings,
		Build:    build.Read(),
	})
	if failure != nil {
		test.Fatal(failure)
	}
	logger.Info("ready", slog.String("outcome", "SUCCESS"))

	var record map[string]any
	if failure := json.Unmarshal(output.Bytes(), &record); failure != nil {
		test.Fatal(failure)
	}
	expected := map[string]any{
		"service": "drizz",
		"level":   "INFO",
		"message": "ready",
		"outcome": "SUCCESS",
	}
	for key, value := range expected {
		if record[key] != value {
			test.Fatalf("%s = %#v", key, record[key])
		}
	}
	for _, key := range []string{"version", "revision"} {
		if record[key] == "" {
			test.Fatalf("%s is empty", key)
		}
	}
	value, valid := record["timestamp"].(string)
	if !valid {
		test.Fatalf("timestamp = %#v", record["timestamp"])
	}
	if _, failure := time.Parse(time.RFC3339Nano, value); failure != nil {
		test.Fatalf("timestamp = %q: %v", value, failure)
	}
	if _, found := record["time"]; found {
		test.Fatal("legacy time key is present")
	}
	if _, found := record["msg"]; found {
		test.Fatal("legacy msg key is present")
	}
}

func TestSource(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	settings, failure := logging.New("debug")
	if failure != nil {
		test.Fatal(failure)
	}
	logger, failure := platform.New(platform.Options{
		Output:   &output,
		Settings: settings,
		Build:    build.Read(),
	})
	if failure != nil {
		test.Fatal(failure)
	}
	logger.Info("ready")

	var record map[string]any
	if failure := json.Unmarshal(output.Bytes(), &record); failure != nil {
		test.Fatal(failure)
	}
	source, valid := record["source"].(map[string]any)
	if !valid {
		test.Fatalf("source = %#v", record["source"])
	}
	for _, key := range []string{"file", "function", "line"} {
		if source[key] == nil || source[key] == "" {
			test.Fatalf("source.%s = %#v", key, source[key])
		}
	}
	value, valid := record["timestamp"].(string)
	if !valid || len(value) == 0 || value[len(value)-1] != 'Z' {
		test.Fatalf("timestamp is not UTC: %#v", record["timestamp"])
	}
}

func TestRedaction(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	settings, failure := logging.New("debug")
	if failure != nil {
		test.Fatal(failure)
	}
	logger, failure := platform.New(platform.Options{
		Output:   &output,
		Settings: settings,
		Build:    build.Read(),
	})
	if failure != nil {
		test.Fatal(failure)
	}
	logger.Info("request", slog.String("access_token", "private"))
	logger.Info("session", slog.String("session_id", "session_123"))
	logger.Info("place", slog.String("capital", "Delhi"))

	decoder := json.NewDecoder(&output)
	var secret map[string]any
	if failure := decoder.Decode(&secret); failure != nil {
		test.Fatal(failure)
	}
	if secret["access_token"] != "[REDACTED]" {
		test.Fatalf("access_token = %v", secret["access_token"])
	}
	var session map[string]any
	if failure := decoder.Decode(&session); failure != nil {
		test.Fatal(failure)
	}
	if session["session_id"] != "session_123" {
		test.Fatalf("session_id = %v", session["session_id"])
	}
	var place map[string]any
	if failure := decoder.Decode(&place); failure != nil {
		test.Fatal(failure)
	}
	if place["capital"] != "Delhi" {
		test.Fatalf("capital = %v", place["capital"])
	}
}

func TestNested(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	settings, failure := logging.New("debug")
	if failure != nil {
		test.Fatal(failure)
	}
	logger, failure := platform.New(platform.Options{
		Output:   &output,
		Settings: settings,
		Build:    build.Read(),
	})
	if failure != nil {
		test.Fatal(failure)
	}
	logger.Info("request", slog.Group("authorization", slog.String("value", "secret"), slog.String("scheme", "bearer")))

	var record map[string]any
	if failure := json.Unmarshal(output.Bytes(), &record); failure != nil {
		test.Fatal(failure)
	}
	group, valid := record["authorization"].(map[string]any)
	if !valid {
		test.Fatalf("authorization = %#v", record["authorization"])
	}
	if group["value"] != "[REDACTED]" || group["scheme"] != "[REDACTED]" {
		test.Fatalf("a sensitive group was not redacted: %#v", group)
	}
}

func TestPolicy(test *testing.T) {
	allocations := testing.AllocsPerRun(1000, func() {
		platform.Policy{}.Handler()(nil, slog.String("outcome", "SUCCESS"))
	})
	if allocations != 0 {
		test.Fatalf("allocations = %v", allocations)
	}
}

func TestOutput(test *testing.T) {
	test.Parallel()

	settings, failure := logging.New("")
	if failure != nil {
		test.Fatal(failure)
	}
	_, failure = platform.New(platform.Options{Settings: settings, Build: build.Read()})
	if failure == nil {
		test.Fatal("missing output was accepted")
	}
}
