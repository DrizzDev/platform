package architecture_test

import (
	"go/ast"
	"path/filepath"
	"strings"
	"testing"
	"unicode"
)

func TestAccess(test *testing.T) {
	test.Parallel()

	repository := repository{root: filepath.Join("..", "..", "internal"), test: test}
	repository.walk(func(source entry) {
		if source.file.IsDir() || filepath.Ext(source.path) != ".go" || strings.HasSuffix(source.path, "_test.go") {
			return
		}
		document := document{path: source.path, file: repository.parse(source.path)}
		if strings.Contains(filepath.ToSlash(source.path), "/domain/") {
			repository.fields(document)
		}
		repository.system(document)
	})
}

type document struct {
	path string
	file *ast.File
}

func (repository repository) fields(document document) {
	ast.Inspect(document.file, func(node ast.Node) bool {
		field, valid := node.(*ast.Field)
		if !valid {
			return true
		}
		for _, name := range field.Names {
			if name.IsExported() {
				repository.test.Errorf("%s exposes mutable domain field %q", document.path, name.Name)
			}
		}
		return true
	})
}

func (repository repository) system(document document) {
	allowed := strings.Contains(filepath.ToSlash(document.path), "/infrastructure/") ||
		strings.Contains(filepath.ToSlash(document.path), "/adapter/") ||
		strings.Contains(filepath.ToSlash(document.path), "/platform/filesystem/")
	ast.Inspect(document.file, func(node ast.Node) bool {
		call, valid := node.(*ast.CallExpr)
		if !valid {
			return true
		}
		selector, valid := call.Fun.(*ast.SelectorExpr)
		if !valid {
			return true
		}
		owner, valid := selector.X.(*ast.Ident)
		if !valid || owner.Name != "os" || allowed {
			return true
		}
		if unicode.IsUpper([]rune(selector.Sel.Name)[0]) {
			repository.test.Errorf("%s accesses os.%s outside an approved boundary", document.path, selector.Sel.Name)
		}
		return true
	})
}
