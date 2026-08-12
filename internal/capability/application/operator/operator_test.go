package operator_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"
	"time"

	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/outcome"
	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/artifact"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/sqlite"
	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/app"
	"github.com/DrizzDev/platform/internal/device/domain/bundle"
	"github.com/DrizzDev/platform/internal/device/domain/capture"
	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/emulator"
	"github.com/DrizzDev/platform/internal/device/domain/format"
	"github.com/DrizzDev/platform/internal/device/domain/geo"
	"github.com/DrizzDev/platform/internal/device/domain/image"
	"github.com/DrizzDev/platform/internal/device/domain/parcel"
	"github.com/DrizzDev/platform/internal/device/domain/pinch"
	"github.com/DrizzDev/platform/internal/device/domain/platform"
	"github.com/DrizzDev/platform/internal/device/domain/press"
	"github.com/DrizzDev/platform/internal/device/domain/serial"
	"github.com/DrizzDev/platform/internal/device/domain/swipe"
	"github.com/DrizzDev/platform/internal/device/domain/text"
	"github.com/DrizzDev/platform/internal/device/domain/touch"
)

// bridge is a fake device that returns a fixed capture and device set, so the operator can be driven through the real
// device control flow without a physical device.
type bridge struct {
	frame   capture.Capture
	devices []device.Device
}

func (bridge bridge) Screenshot(context.Context, device.Device) (capture.Capture, error) {
	return bridge.frame, nil
}

func (bridge bridge) Snapshot(context.Context, device.Device) (capture.Capture, string, error) {
	return bridge.frame, "<hierarchy/>", nil
}

func (bridge bridge) Hierarchy(context.Context, device.Device) (string, error) {
	return "<hierarchy/>", nil
}

func (bridge bridge) Dimensions(context.Context, device.Device) (int, int, error) {
	return 1080, 2400, nil
}

func (bridge bridge) Tap(context.Context, touch.Touch) error { return nil }

func (bridge bridge) Swipe(context.Context, swipe.Swipe) error { return nil }

func (bridge bridge) Pinch(context.Context, pinch.Pinch) error { return nil }

func (bridge bridge) Press(context.Context, press.Press) error { return nil }

func (bridge bridge) Type(context.Context, text.Text) error { return nil }

func (bridge bridge) Clear(context.Context, device.Device) error { return nil }

func (bridge bridge) Back(context.Context, device.Device) error { return nil }

func (bridge bridge) Home(context.Context, device.Device) error { return nil }

func (bridge bridge) Locate(context.Context, geo.Fix) error { return nil }

func (bridge bridge) Install(context.Context, parcel.Parcel) error { return nil }

func (bridge bridge) Launch(context.Context, bundle.Bundle) error { return nil }

func (bridge bridge) Terminate(context.Context, bundle.Bundle) error { return nil }

func (bridge bridge) Wipe(context.Context, bundle.Bundle) error { return nil }

func (bridge bridge) Installed(context.Context, device.Device) ([]app.App, error) {
	return []app.App{app.New(app.Input{Id: "com.example", Name: "Example", Note: "1.0"})}, nil
}

func (bridge bridge) Running(context.Context, device.Device) ([]app.App, error) {
	return []app.App{app.New(app.Input{Id: "com.example", Name: "Example", Note: "123"})}, nil
}

func (bridge bridge) Foreground(context.Context, device.Device) (string, error) {
	return "com.example", nil
}

func (bridge bridge) Url(context.Context, device.Device) (string, error) {
	return "https://example.com", nil
}

func (bridge bridge) Disk(context.Context, device.Device) (int, error) { return 4096, nil }

func (bridge bridge) Images(context.Context, platform.Platform) ([]string, error) {
	return []string{"Pixel_7_API_34"}, nil
}

func (bridge bridge) Boot(context.Context, emulator.Boot) error { return nil }

func (bridge bridge) Pause(context.Context, device.Device) error { return nil }

func (bridge bridge) Resume(context.Context, device.Device) error { return nil }

func (bridge bridge) List(context.Context) ([]device.Device, error) {
	return bridge.devices, nil
}

type clock struct{}

func (clock) Now() time.Time { return time.Unix(1_000_000, 0) }

type kit struct {
	test *testing.T
}

func (kit kit) device(id string) device.Device {
	kit.test.Helper()
	address, failure := serial.New(id)
	if failure != nil {
		kit.test.Fatal(failure)
	}
	made, failure := device.New(device.Input{Serial: address, Platform: platform.Android})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return made
}

func (kit kit) capture(data string) capture.Capture {
	kit.test.Helper()
	picture, failure := image.New([]byte(data))
	if failure != nil {
		kit.test.Fatal(failure)
	}
	made, failure := capture.New(capture.Input{Image: picture, Format: format.Png, Width: 1080, Height: 2400})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return made
}

func (kit kit) build(desk bridge) (operator.Operator, sqlite.Store, artifact.Store) {
	kit.test.Helper()
	dir := kit.test.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	tracer := tracenoop.NewTracerProvider().Tracer("test")
	meter := metricnoop.NewMeterProvider().Meter("test")

	store, failure := sqlite.New(context.Background(), sqlite.Options{
		Path: filepath.Join(dir, "capture.db"), Logger: logger, Tracer: tracer, Meter: meter,
	})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	kit.test.Cleanup(func() { _ = store.Close() })
	vault, failure := artifact.New(artifact.Options{Root: dir, Logger: logger, Tracer: tracer, Meter: meter})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	recorder, failure := recording.New(recording.Options{
		Writer: store, Sink: vault, Keeper: store, Clock: clock{}, Logger: logger, Lease: time.Minute,
	})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	flow, failure := control.New(control.Options{Bridge: desk})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	performed, failure := operator.New(operator.Options{Flow: flow, Recorder: recorder, Logger: logger})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return performed, store, vault
}

func TestScreenshot(test *testing.T) {
	test.Parallel()

	kit := kit{test: test}
	desk := bridge{frame: kit.capture("image-bytes"), devices: []device.Device{kit.device("s-1")}}
	performed, store, vault := kit.build(desk)
	scope := context.Background()

	shot, failure := performed.Screenshot(scope, operator.Target{Serial: "s-1"})
	if failure != nil {
		test.Fatal(failure)
	}
	if string(shot.Image) != "image-bytes" || shot.Format != "PNG" {
		test.Fatalf("shot = %q %q", shot.Image, shot.Format)
	}

	recorded, failure := store.Pending(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(recorded) != 1 || recorded[0].Category() != category.Screen {
		test.Fatalf("recorded %d entries, want one screen capture", len(recorded))
	}
	blob, failure := vault.Get(scope, recorded[0].Artifact())
	if failure != nil {
		test.Fatal(failure)
	}
	if string(blob) != "image-bytes" {
		test.Fatalf("recorded artifact = %q", blob)
	}
}

func TestTap(test *testing.T) {
	test.Parallel()

	kit := kit{test: test}
	performed, store, _ := kit.build(bridge{devices: []device.Device{kit.device("s-1")}})
	scope := context.Background()

	if _, failure := performed.Tap(scope, operator.Contact{Serial: "s-1", X: 100, Y: 200}); failure != nil {
		test.Fatal(failure)
	}

	recorded, failure := store.Pending(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(recorded) != 1 || recorded[0].Category() != category.Tool {
		test.Fatalf("recorded %d entries, want one tool action", len(recorded))
	}
}

func TestScreenshotUnknownDevice(test *testing.T) {
	test.Parallel()

	kit := kit{test: test}
	performed, _, _ := kit.build(bridge{devices: []device.Device{kit.device("s-1")}})

	_, failure := performed.Screenshot(context.Background(), operator.Target{Serial: "s-2"})
	var refusal operator.Refusal
	if !errors.As(failure, &refusal) || refusal.Code != outcome.Missing {
		test.Fatalf("unknown device = %v, want a not-found refusal", failure)
	}
}

// broken is a recorder whose record cannot be opened, so the operator's drop path can be exercised.
type broken struct{}

func (broken) Begin() (*recording.Execution, error) {
	return nil, errors.New("recorder unavailable")
}

func TestRecordDropObservable(test *testing.T) {
	test.Parallel()

	kit := kit{test: test}
	var log bytes.Buffer
	desk := bridge{frame: kit.capture("image-bytes"), devices: []device.Device{kit.device("s-1")}}
	flow, failure := control.New(control.Options{Bridge: desk})
	if failure != nil {
		test.Fatal(failure)
	}
	performed, failure := operator.New(operator.Options{
		Flow: flow, Recorder: broken{}, Logger: slog.New(slog.NewJSONHandler(&log, nil)),
	})
	if failure != nil {
		test.Fatal(failure)
	}

	shot, failure := performed.Screenshot(context.Background(), operator.Target{Serial: "s-1"})
	if failure != nil {
		test.Fatalf("a failed record must not fail the capture: %v", failure)
	}
	if string(shot.Image) != "image-bytes" {
		test.Fatalf("shot = %q", shot.Image)
	}
	if !strings.Contains(log.String(), "capability.record.dropped") {
		test.Fatal("a dropped record must be observable, not silent")
	}
}

func TestDevices(test *testing.T) {
	test.Parallel()

	kit := kit{test: test}
	performed, _, _ := kit.build(bridge{devices: []device.Device{kit.device("s-1"), kit.device("s-2")}})

	roster, failure := performed.Devices(context.Background())
	if failure != nil {
		test.Fatal(failure)
	}
	if len(roster.Serials) != 2 {
		test.Fatalf("roster = %v, want two serials", roster.Serials)
	}
}
