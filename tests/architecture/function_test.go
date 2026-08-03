package architecture_test

import (
	"go/ast"
	"path/filepath"
	"strings"
	"testing"
)

func TestFunctions(test *testing.T) {
	test.Parallel()

	repository := repository{root: filepath.Join("..", ".."), test: test}
	repository.walk(func(source entry) {
		if source.file.IsDir() || filepath.Ext(source.path) != ".go" {
			return
		}
		file := repository.parse(source.path)
		for _, declaration := range file.Decls {
			function, valid := declaration.(*ast.FuncDecl)
			if !valid {
				continue
			}
			repository.check(inspection{path: source.path, function: function})
		}
	})
}

type inspection struct {
	path     string
	function *ast.FuncDecl
}

func (repository repository) check(inspection inspection) {
	count := repository.parameters(inspection.function.Type.Params)
	if count > 2 {
		repository.test.Errorf(
			"%s: %s accepts more than two parameters",
			inspection.path,
			inspection.function.Name.Name,
		)
	}
	if count == 2 && !repository.contextual(inspection.function.Type.Params) {
		repository.test.Errorf(
			"%s: %s must replace positional parameters with a typed input",
			inspection.path,
			inspection.function.Name.Name,
		)
	}
	seam := strings.HasSuffix(inspection.path, "export_test.go")
	if inspection.function.Recv == nil && !seam && !repository.entry(inspection.function.Name.Name) {
		repository.test.Errorf(
			"%s: package function %s is not a constructor, test, or entry point",
			inspection.path,
			inspection.function.Name.Name,
		)
	}
}

func (repository repository) contextual(fields *ast.FieldList) bool {
	if fields == nil || len(fields.List) == 0 {
		return false
	}
	selector, valid := fields.List[0].Type.(*ast.SelectorExpr)
	if !valid || selector.Sel.Name != "Context" {
		return false
	}
	owner, valid := selector.X.(*ast.Ident)
	return valid && owner.Name == "context"
}

func (repository repository) parameters(fields *ast.FieldList) int {
	if fields == nil {
		return 0
	}
	count := 0
	for _, field := range fields.List {
		if len(field.Names) == 0 {
			count++
		} else {
			count += len(field.Names)
		}
	}
	return count
}

func (repository repository) entry(name string) bool {
	return name == "New" ||
		name == "Read" ||
		name == "Serve" ||
		name == "main" ||
		name == "run" ||
		strings.HasPrefix(name, "Test") ||
		strings.HasPrefix(name, "Fuzz") ||
		strings.HasPrefix(name, "Benchmark")
}
