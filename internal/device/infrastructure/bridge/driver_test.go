package bridge_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"sync/atomic"
	"testing"

	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/platform"
	"github.com/DrizzDev/platform/internal/device/domain/point"
	"github.com/DrizzDev/platform/internal/device/domain/serial"
	"github.com/DrizzDev/platform/internal/device/domain/touch"
	"github.com/DrizzDev/platform/internal/device/infrastructure/bridge"
)

type maker struct {
	test *testing.T
}

func (maker maker) target() device.Device {
	maker.test.Helper()
	handle, failure := serial.New("emulator-5554")
	if failure != nil {
		maker.test.Fatal(failure)
	}
	target, failure := device.New(device.Input{Platform: platform.Android, Serial: handle})
	if failure != nil {
		maker.test.Fatal(failure)
	}
	return target
}

func (maker maker) contact() touch.Touch {
	maker.test.Helper()
	spot, failure := point.New(point.Input{X: 540, Y: 1200})
	if failure != nil {
		maker.test.Fatal(failure)
	}
	contact, failure := touch.New(touch.Input{Device: maker.target(), Point: spot})
	if failure != nil {
		maker.test.Fatal(failure)
	}
	return contact
}

func (fixture fixture) driver() *bridge.Driver {
	fixture.test.Helper()
	return bridge.Wrap(fixture.channel())
}

func TestList(test *testing.T) {
	test.Parallel()

	driver := fixture{test: test, handler: func(request call) response {
		return response{Version: "2.0", ID: request.ID, Result: json.RawMessage(`[{"platform":"android","udid":"emulator-5554"}]`)}
	}}.driver()
	defer func() { _ = driver.Close() }()

	devices, failure := driver.List(context.Background())
	if failure != nil {
		test.Fatal(failure)
	}
	if len(devices) != 1 || devices[0].Serial().String() != "emulator-5554" || devices[0].Platform() != platform.Android {
		test.Fatalf("devices = %v", devices)
	}
}

func TestScreenshot(test *testing.T) {
	test.Parallel()

	data := base64.StdEncoding.EncodeToString([]byte{1, 2, 3})
	driver := fixture{test: test, handler: func(request call) response {
		body, _ := json.Marshal(map[string]any{"format": "png", "width": 1080, "height": 2400, "data": data})
		return response{Version: "2.0", ID: request.ID, Result: body}
	}}.driver()
	defer func() { _ = driver.Close() }()

	shot, failure := driver.Screenshot(context.Background(), maker{test: test}.target())
	if failure != nil {
		test.Fatal(failure)
	}
	if shot.Width() != 1080 || shot.Height() != 2400 || shot.Image().Size() != 3 {
		test.Fatalf("capture = %dx%d/%d", shot.Width(), shot.Height(), shot.Image().Size())
	}
}

func TestTap(test *testing.T) {
	test.Parallel()

	driver := fixture{test: test, handler: func(request call) response {
		return response{Version: "2.0", ID: request.ID, Result: json.RawMessage(`{}`)}
	}}.driver()
	defer func() { _ = driver.Close() }()

	if failure := driver.Tap(context.Background(), maker{test: test}.contact()); failure != nil {
		test.Fatal(failure)
	}
}

func TestKind(test *testing.T) {
	test.Parallel()

	driver := fixture{test: test, handler: func(request call) response {
		return response{Version: "2.0", ID: request.ID,
			Error: &defect{Code: -32000, Message: "adb secure surface", Data: carrier{Kind: "SecureScreenCaptureError"}}}
	}}.driver()
	defer func() { _ = driver.Close() }()

	_, failure := driver.Screenshot(context.Background(), maker{test: test}.target())
	if !errors.Is(failure, control.Protected{}) {
		test.Fatalf("failure = %v", failure)
	}
}

func TestUnsupported(test *testing.T) {
	test.Parallel()

	driver := fixture{test: test, handler: func(request call) response {
		return response{Version: "2.0", ID: request.ID,
			Error: &defect{Code: -32000, Message: "pressButton is not supported", Data: carrier{Kind: "UnsupportedOperationError"}}}
	}}.driver()
	defer func() { _ = driver.Close() }()

	_, failure := driver.Screenshot(context.Background(), maker{test: test}.target())
	if !errors.Is(failure, control.Incompatible{}) {
		test.Fatalf("an unsupported operation must map to the incompatible outcome, got %v", failure)
	}
}

func TestMalformed(test *testing.T) {
	test.Parallel()

	driver := fixture{test: test, handler: func(request call) response {
		return response{Version: "2.0", ID: request.ID, Result: json.RawMessage(`{"format":"gif","width":10,"height":10,"data":"Zm9v"}`)}
	}}.driver()
	defer func() { _ = driver.Close() }()

	_, failure := driver.Screenshot(context.Background(), maker{test: test}.target())
	if !errors.Is(failure, control.Failed{}) {
		test.Fatalf("failure = %v", failure)
	}
}

func TestReadRetry(test *testing.T) {
	test.Parallel()

	var builds atomic.Int64
	channel := bridge.Attach(func(_ context.Context) (io.Writer, io.Reader, func() error) {
		first := builds.Add(1) == 1
		reader, writer := io.Pipe()
		outward, inward := io.Pipe()
		stage := &puppet{output: inward, crash: first, handler: func(request call) response {
			return response{Version: "2.0", ID: request.ID, Result: json.RawMessage(`[{"platform":"android","udid":"emulator-5554"}]`)}
		}}
		go stage.serve(reader)
		return writer, outward, func() error {
			_ = writer.Close()
			_ = inward.Close()
			return nil
		}
	})
	driver := bridge.Wrap(channel)
	defer func() { _ = driver.Close() }()

	devices, failure := driver.List(context.Background())
	if failure != nil {
		test.Fatalf("retry failed: %v", failure)
	}
	if len(devices) != 1 || builds.Load() < 2 {
		test.Fatalf("devices=%d builds=%d", len(devices), builds.Load())
	}
}

func TestTapNoRetry(test *testing.T) {
	test.Parallel()

	var builds atomic.Int64
	channel := bridge.Attach(func(_ context.Context) (io.Writer, io.Reader, func() error) {
		builds.Add(1)
		reader, writer := io.Pipe()
		outward, inward := io.Pipe()
		stage := &puppet{output: inward, crash: true, handler: echo}
		go stage.serve(reader)
		return writer, outward, func() error {
			_ = writer.Close()
			_ = inward.Close()
			return nil
		}
	})
	driver := bridge.Wrap(channel)
	defer func() { _ = driver.Close() }()

	failure := driver.Tap(context.Background(), maker{test: test}.contact())
	if !errors.Is(failure, control.Unavailable{}) {
		test.Fatalf("failure = %v", failure)
	}
	if builds.Load() != 1 {
		test.Fatalf("tap retried: %d builds", builds.Load())
	}
}
