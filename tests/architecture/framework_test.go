package architecture_test

import (
	"go/ast"
	"path/filepath"
	"strings"
	"testing"
)

type confinement struct {
	framework string
	owner     string
}

type usage struct {
	path string
	line *ast.ImportSpec
}

var frameworks = []confinement{
	{framework: "github.com/spf13/cobra", owner: "/internal/transport/cli"},
	{framework: "github.com/modelcontextprotocol/go-sdk", owner: "/internal/transport/mcp"},
	{framework: "github.com/getsentry/", owner: "/internal/platform/reporting/sentry"},
	{framework: "github.com/samber/slog-sentry", owner: "/internal/platform/reporting/sentry"},
	{framework: "go.opentelemetry.io/otel/sdk", owner: "/internal/platform/telemetry"},
	{framework: "go.opentelemetry.io/otel/exporters", owner: "/internal/platform/telemetry"},
}

func TestFramework(test *testing.T) {
	test.Parallel()

	repository := repository{root: filepath.Join("..", "..", "internal"), test: test}
	repository.walk(func(source entry) {
		if source.file.IsDir() || filepath.Ext(source.path) != ".go" || strings.HasSuffix(source.path, "_test.go") {
			return
		}
		path := filepath.ToSlash(source.path)
		for _, line := range repository.parse(source.path).Imports {
			repository.confine(usage{path: path, line: line})
		}
	})
}

func (repository repository) confine(usage usage) {
	value := strings.Trim(usage.line.Path.Value, `"`)
	for _, rule := range frameworks {
		if strings.HasPrefix(value, rule.framework) && !strings.Contains(usage.path, rule.owner) {
			repository.test.Errorf("%s imports %s outside %s", usage.path, value, rule.owner)
		}
	}
	if value == "os" && usage.line.Name != nil {
		repository.test.Errorf("%s imports os under an alias", usage.path)
	}
}
