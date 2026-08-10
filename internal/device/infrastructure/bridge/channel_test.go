package bridge_test

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/device/infrastructure/bridge"
)

// call, response, and defect are the test's own view of the wire protocol, so the
// fake speaks the real JSON contract independently of the package internals.
type call struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	ID     int64           `json:"id"`
}

type response struct {
	Version string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *defect         `json:"error,omitempty"`
	ID      int64           `json:"id"`
}

type defect struct {
	Message string  `json:"message"`
	Data    carrier `json:"data"`
	Code    int     `json:"code"`
}

type carrier struct {
	Kind string `json:"kind"`
}

// puppet is a scripted fake sidecar. It answers health, runs each call handler in
// its own goroutine so the read loop stays responsive to `$/cancel`, and can
// close its output to simulate a crash.
type puppet struct {
	output  *io.PipeWriter
	handler func(call) response
	cancels chan int64
	mutex   sync.Mutex
}

func (puppet *puppet) serve(requests io.Reader) {
	reader := bufio.NewReader(requests)
	for {
		line, failure := reader.ReadBytes('\n')
		if len(line) > 1 {
			puppet.dispatch(line)
		}
		if failure != nil {
			return
		}
	}
}

func (puppet *puppet) dispatch(line []byte) {
	var request call
	if json.Unmarshal(line, &request) != nil {
		return
	}
	switch request.Method {
	case "$/cancel":
		var target struct {
			ID int64 `json:"id"`
		}
		_ = json.Unmarshal(request.Params, &target)
		select {
		case puppet.cancels <- target.ID:
		default:
		}
	case "crash":
		puppet.mutex.Lock()
		_ = puppet.output.Close()
		puppet.mutex.Unlock()
	default:
		go puppet.reply(request)
	}
}

func (puppet *puppet) reply(request call) {
	var answer response
	if request.Method == "health" {
		answer = response{Version: "2.0", ID: request.ID, Result: json.RawMessage(`{"status":"ok","version":"t","protocol":1}`)}
	} else {
		answer = puppet.handler(request)
	}
	line, _ := json.Marshal(answer)
	puppet.mutex.Lock()
	_, _ = puppet.output.Write(append(line, '\n'))
	puppet.mutex.Unlock()
}

type fixture struct {
	test    *testing.T
	handler func(call) response
	cancels chan int64
}

func (fixture fixture) channel() *bridge.Channel {
	fixture.test.Helper()
	return bridge.Attach(func(_ context.Context) (io.Writer, io.Reader, func() error) {
		reader, writer := io.Pipe()
		outward, inward := io.Pipe()
		stage := &puppet{output: inward, handler: fixture.handler, cancels: fixture.cancels}
		go stage.serve(reader)
		return writer, outward, func() error {
			_ = writer.Close()
			_ = inward.Close()
			return nil
		}
	})
}

var echo = func(request call) response {
	return response{Version: "2.0", ID: request.ID, Result: request.Params}
}

func TestInvoke(test *testing.T) {
	test.Parallel()

	channel := fixture{test: test, handler: echo}.channel()
	defer func() { _ = channel.Close() }()

	result, failure := channel.Invoke(context.Background(),
		bridge.Request{Method: "device.tap", Params: map[string]int{"x": 5}})
	if failure != nil {
		test.Fatal(failure)
	}
	if string(result.Result) != `{"x":5}` {
		test.Fatalf("result = %s", result.Result)
	}
}

func TestFault(test *testing.T) {
	test.Parallel()

	channel := fixture{test: test, handler: func(request call) response {
		return response{Version: "2.0", ID: request.ID,
			Error: &defect{Code: -32000, Message: "internal adb path leaked here", Data: carrier{Kind: "DeviceNotFoundError"}}}
	}}.channel()
	defer func() { _ = channel.Close() }()

	_, failure := channel.Invoke(context.Background(), bridge.Request{Method: "device.screenshot"})
	var fault bridge.Fault
	if !errors.As(failure, &fault) {
		test.Fatalf("failure = %v", failure)
	}
	if fault.Kind() != "DeviceNotFoundError" {
		test.Fatalf("kind = %q", fault.Kind())
	}
}

func TestCancel(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test, cancels: make(chan int64, 4), handler: func(request call) response {
		time.Sleep(300 * time.Millisecond)
		return echo(request)
	}}
	channel := kit.channel()
	defer func() { _ = channel.Close() }()

	scope, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	_, failure := channel.Invoke(scope, bridge.Request{Method: "device.screenshot"})
	if !errors.Is(failure, context.Canceled) {
		test.Fatalf("failure = %v", failure)
	}
	select {
	case <-kit.cancels:
	case <-time.After(time.Second):
		test.Fatal("the sidecar was never told to cancel")
	}
}

func TestRestart(test *testing.T) {
	test.Parallel()

	channel := fixture{test: test, handler: echo}.channel()
	defer func() { _ = channel.Close() }()

	_, failure := channel.Invoke(context.Background(), bridge.Request{Method: "crash"})
	if !errors.Is(failure, bridge.Unavailable{}) {
		test.Fatalf("crash failure = %v", failure)
	}
	result, failure := channel.Invoke(context.Background(),
		bridge.Request{Method: "device.tap", Params: map[string]int{"y": 9}})
	if failure != nil {
		test.Fatalf("restart failed: %v", failure)
	}
	if string(result.Result) != `{"y":9}` {
		test.Fatalf("result = %s", result.Result)
	}
}

func TestOversize(test *testing.T) {
	test.Parallel()

	line := `"` + strings.Repeat("x", 512) + `"` + "\n"
	_, failure := bridge.Scan(bridge.Probe{Source: strings.NewReader(line), Bound: 256})
	if !errors.Is(failure, bridge.Oversize{}) {
		test.Fatalf("failure = %v", failure)
	}
}

func TestRecycle(test *testing.T) {
	test.Parallel()

	var builds atomic.Int64
	wedge := func(request call) response {
		time.Sleep(500 * time.Millisecond)
		return echo(request)
	}
	channel := bridge.Attach(func(_ context.Context) (io.Writer, io.Reader, func() error) {
		builds.Add(1)
		reader, writer := io.Pipe()
		outward, inward := io.Pipe()
		stage := &puppet{output: inward, handler: wedge}
		go stage.serve(reader)
		return writer, outward, func() error {
			_ = writer.Close()
			_ = inward.Close()
			return nil
		}
	})
	defer func() { _ = channel.Close() }()

	for attempt := 0; attempt <= bridge.Patience; attempt++ {
		scope, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, _ = channel.Invoke(scope, bridge.Request{Method: "device.screenshot"})
		cancel()
	}
	if builds.Load() < 2 {
		test.Fatalf("a wedged sidecar was not recycled: %d builds", builds.Load())
	}
}

func TestClose(test *testing.T) {
	test.Parallel()

	channel := fixture{test: test, handler: echo}.channel()
	if _, failure := channel.Invoke(context.Background(), bridge.Request{Method: "device.tap"}); failure != nil {
		test.Fatal(failure)
	}
	if failure := channel.Close(); failure != nil {
		test.Fatal(failure)
	}
	_, failure := channel.Invoke(context.Background(), bridge.Request{Method: "device.tap"})
	if !errors.Is(failure, bridge.Closed{}) {
		test.Fatalf("failure = %v", failure)
	}
}
