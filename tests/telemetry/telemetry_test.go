package telemetry_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/platform/build"
	configuration "github.com/DrizzDev/platform/internal/platform/configuration/telemetry"
	"github.com/DrizzDev/platform/internal/platform/telemetry"
)

const deadline = 15 * time.Second

func TestExport(test *testing.T) {
	test.Parallel()

	collector := &collector{}
	server := httptest.NewServer(collector.handler())
	defer server.Close()

	scope, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	provider := harness{test: test}.provider(scope, server.URL)
	(harness{test: test}).emit(scope, provider)
	if failure := provider.Close(scope); failure != nil {
		test.Fatalf("successful export reported a failure: %v", failure)
	}
	traces, metrics := collector.totals()
	if traces == 0 || metrics == 0 {
		test.Fatalf("collector received traces=%d metrics=%d", traces, metrics)
	}
}

func TestDelivery(test *testing.T) {
	test.Parallel()

	collector := &collector{status: http.StatusBadRequest}
	server := httptest.NewServer(collector.handler())
	defer server.Close()

	scope, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	provider := harness{test: test}.provider(scope, server.URL)
	(harness{test: test}).emit(scope, provider)
	if failure := provider.Close(scope); failure == nil {
		test.Fatal("a rejected export did not surface a failure")
	}
}

func TestCancellation(test *testing.T) {
	test.Parallel()

	server := httptest.NewServer((&collector{}).handler())
	defer server.Close()

	scope, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	provider := harness{test: test}.provider(scope, server.URL)
	(harness{test: test}).emit(scope, provider)

	stopped, halt := context.WithCancel(context.Background())
	halt()
	if failure := provider.Close(stopped); failure == nil {
		test.Fatal("cancelled shutdown did not report the cancellation")
	}
}

func TestRepeated(test *testing.T) {
	test.Parallel()

	server := httptest.NewServer((&collector{}).handler())
	defer server.Close()

	scope, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	provider := harness{test: test}.provider(scope, server.URL)
	if failure := provider.Close(scope); failure != nil {
		test.Fatalf("first shutdown: %v", failure)
	}
	if failure := provider.Close(scope); failure != nil {
		test.Fatalf("repeated shutdown must remain idempotent: %v", failure)
	}
}

func TestSilent(test *testing.T) {
	test.Parallel()

	collector := &collector{}
	server := httptest.NewServer(collector.handler())
	defer server.Close()

	scope, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	settings, failure := configuration.New(configuration.Input{})
	if failure != nil {
		test.Fatal(failure)
	}
	provider, failure := telemetry.New(scope, telemetry.Options{Settings: settings, Build: build.Read()})
	if failure != nil {
		test.Fatal(failure)
	}
	(harness{test: test}).emit(scope, provider)
	if failure := provider.Close(scope); failure != nil {
		test.Fatal(failure)
	}
	if traces, metrics := collector.totals(); traces != 0 || metrics != 0 {
		test.Fatalf("disabled telemetry performed network work: traces=%d metrics=%d", traces, metrics)
	}
}

type harness struct {
	test *testing.T
}

func (harness harness) emit(scope context.Context, provider telemetry.Provider) {
	_, span := provider.Tracer().Start(scope, "operation")
	span.End()
	counter, failure := provider.Meter().Int64Counter("operation")
	if failure == nil {
		counter.Add(scope, 1)
	}
}

func (harness harness) provider(scope context.Context, endpoint string) telemetry.Provider {
	harness.test.Helper()
	settings, failure := configuration.New(configuration.Input{Exporter: "OTLP", Endpoint: endpoint})
	if failure != nil {
		harness.test.Fatal(failure)
	}
	provider, failure := telemetry.New(scope, telemetry.Options{Settings: settings, Build: build.Read()})
	if failure != nil {
		harness.test.Fatal(failure)
	}
	return provider
}

type collector struct {
	mutex   sync.Mutex
	traces  int
	metrics int
	status  int
}

func (collector *collector) handler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		body, _ := io.ReadAll(request.Body)
		collector.mutex.Lock()
		switch {
		case strings.Contains(request.URL.Path, "traces") && len(body) != 0:
			collector.traces++
		case strings.Contains(request.URL.Path, "metrics") && len(body) != 0:
			collector.metrics++
		}
		status := collector.status
		collector.mutex.Unlock()
		if status == 0 {
			status = http.StatusOK
		}
		writer.WriteHeader(status)
	})
}

func (collector *collector) totals() (int, int) {
	collector.mutex.Lock()
	defer collector.mutex.Unlock()
	return collector.traces, collector.metrics
}
