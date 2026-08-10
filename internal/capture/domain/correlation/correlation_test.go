package correlation_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/capture/domain/correlation"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
	"github.com/DrizzDev/platform/internal/capture/domain/trace"
)

type fixture struct {
	test *testing.T
}

func (fixture fixture) trace() trace.Trace {
	fixture.test.Helper()
	value, failure := trace.New("01HEXECUTION")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return value
}

func (fixture fixture) span(value string) span.Span {
	fixture.test.Helper()
	hop, failure := span.New(value)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return hop
}

func TestCorrelationRoot(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	root, failure := correlation.New(correlation.Input{Trace: kit.trace(), Span: kit.span("hop-1"), Mark: mark.Exact})
	if failure != nil {
		test.Fatal(failure)
	}
	if !root.Root() {
		test.Fatal("a parentless hop is not reported as the root")
	}
}

func TestCorrelationChild(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	child, failure := correlation.New(correlation.Input{
		Trace:    kit.trace(),
		Span:     kit.span("hop-2"),
		Parent:   kit.span("hop-1"),
		Sequence: 1,
		Mark:     mark.Inferred,
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if child.Root() || child.Parent().String() != "hop-1" || child.Mark() != mark.Inferred {
		test.Fatalf("child correlation = root:%v parent:%q", child.Root(), child.Parent().String())
	}
}

func TestCorrelationRejects(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	rejected := map[string]correlation.Input{
		"trace":    {Span: kit.span("hop-1"), Mark: mark.Exact},
		"span":     {Trace: kit.trace(), Mark: mark.Exact},
		"sequence": {Trace: kit.trace(), Span: kit.span("hop-1"), Sequence: -1, Mark: mark.Exact},
		"mark":     {Trace: kit.trace(), Span: kit.span("hop-1"), Mark: mark.Mark("MAYBE")},
		"self":     {Trace: kit.trace(), Span: kit.span("hop-1"), Parent: kit.span("hop-1"), Mark: mark.Exact},
	}
	for name, input := range rejected {
		if _, failure := correlation.New(input); failure == nil {
			test.Fatalf("%s correlation was accepted", name)
		}
	}
}
