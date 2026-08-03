package architecture_test

import (
	"go/ast"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsolation(test *testing.T) {
	test.Parallel()

	repository := repository{root: filepath.Join("..", ".."), test: test}
	packages := make(map[string]string)
	repository.walk(func(source entry) {
		if source.file.IsDir() || filepath.Ext(source.path) != ".go" || strings.HasSuffix(source.path, "_test.go") {
			return
		}
		packages[filepath.Dir(source.path)] = repository.parse(source.path).Name.Name
	})
	repository.walk(func(source entry) {
		if source.file.IsDir() || !strings.HasSuffix(source.path, "_test.go") {
			return
		}
		production, found := packages[filepath.Dir(source.path)]
		if !found || repository.parse(source.path).Name.Name != production {
			return
		}
		if filepath.Base(source.path) != "export_test.go" || repository.contains(source.path) {
			test.Errorf("%s must use an external test package", source.path)
		}
	})
}

func (repository repository) contains(path string) bool {
	for _, declaration := range repository.parse(path).Decls {
		function, valid := declaration.(*ast.FuncDecl)
		if valid && (strings.HasPrefix(function.Name.Name, "Test") ||
			strings.HasPrefix(function.Name.Name, "Fuzz") ||
			strings.HasPrefix(function.Name.Name, "Benchmark")) {
			return true
		}
	}
	return false
}
