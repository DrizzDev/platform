package architecture_test

import (
	"go/ast"
	"go/token"
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

// fields enforces value-object immutability: a domain struct must not export a
// mutable field. The typed construction input is exempt — it is a transient DTO
// consumed by New, not retained domain state.
func (repository repository) fields(document document) {
	for _, declaration := range document.file.Decls {
		group, valid := declaration.(*ast.GenDecl)
		if !valid || group.Tok != token.TYPE {
			continue
		}
		for _, specification := range group.Specs {
			repository.immutable(member{document: document, specification: specification})
		}
	}
}

type member struct {
	document      document
	specification ast.Spec
}

func (repository repository) immutable(member member) {
	definition, valid := member.specification.(*ast.TypeSpec)
	if !valid || definition.Name.Name == "Input" {
		return
	}
	structure, valid := definition.Type.(*ast.StructType)
	if !valid {
		return
	}
	for _, field := range structure.Fields.List {
		for _, name := range field.Names {
			if name.IsExported() {
				repository.test.Errorf("%s exposes mutable domain field %q", member.document.path, name.Name)
			}
		}
	}
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
