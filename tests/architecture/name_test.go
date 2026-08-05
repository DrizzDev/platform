package architecture_test

import (
	"go/ast"
	"path/filepath"
	"strings"
	"testing"
	"unicode"
)

func TestNames(test *testing.T) {
	test.Parallel()

	repository := repository{root: filepath.Join("..", ".."), test: test}
	repository.walk(func(source entry) {
		name := source.file.Name()
		if source.file.IsDir() {
			if strings.ContainsAny(name, "_- ") && !repository.integration(source.path) {
				test.Errorf("directory name must be one word: %s", source.path)
			}
			return
		}
		if filepath.Ext(name) != ".go" {
			return
		}
		base := strings.TrimSuffix(strings.TrimSuffix(name, ".go"), "_test")
		if strings.ContainsAny(base, "_- ") {
			test.Errorf("source name must be one word: %s", source.path)
		}
		file := repository.parse(source.path)
		for _, name := range repository.identifiers(file) {
			if repository.prohibited(name) {
				test.Errorf("%s uses prohibited identifier %q", source.path, name)
			}
		}
	})
}

func (repository repository) integration(path string) bool {
	return filepath.Base(path) == "PULL_REQUEST_TEMPLATE" &&
		filepath.Base(filepath.Dir(path)) == ".github"
}

func (repository repository) identifiers(file *ast.File) []string {
	names := &names{}
	ast.Walk(names, file)
	return names.values
}

type names struct {
	values []string
}

func (names *names) Visit(node ast.Node) ast.Visitor {
	switch value := node.(type) {
	case *ast.FuncDecl:
		names.add(value.Name)
	case *ast.TypeSpec:
		names.add(value.Name)
	case *ast.ValueSpec:
		names.list(value.Names)
	case *ast.Field:
		names.list(value.Names)
	case *ast.AssignStmt:
		names.expressions(value.Lhs)
	case *ast.RangeStmt:
		names.expression(value.Key)
		names.expression(value.Value)
	}
	return names
}

func (names *names) list(values []*ast.Ident) {
	for _, value := range values {
		names.add(value)
	}
}

func (names *names) expressions(values []ast.Expr) {
	for _, value := range values {
		names.expression(value)
	}
}

func (names *names) expression(value ast.Expr) {
	if name, valid := value.(*ast.Ident); valid {
		names.add(name)
	}
}

func (names *names) add(value *ast.Ident) {
	if value != nil {
		names.values = append(names.values, value.Name)
	}
}

func (repository repository) prohibited(name string) bool {
	banned := map[string]bool{
		"cfg": true, "ctx": true, "err": true, "helper": true, "impl": true,
		"manager": true, "ok": true, "opts": true, "service": true,
		"util": true, "utils": true,
	}
	if banned[name] {
		return true
	}
	if name == "" || strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "Fuzz") || strings.HasPrefix(name, "Benchmark") {
		return false
	}
	mandated := map[string]bool{"SessionID": true, "RoundTrip": true}
	if mandated[name] {
		return false
	}
	runes := []rune(name)
	for index, value := range runes[1:] {
		if unicode.IsUpper(value) && !unicode.IsUpper(runes[index]) {
			return true
		}
	}
	return false
}
