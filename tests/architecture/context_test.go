package architecture_test

import (
	"go/ast"
	"path/filepath"
	"strings"
	"testing"
)

func TestContext(test *testing.T) {
	test.Parallel()

	repository := repository{root: filepath.Join("..", "..", "internal"), test: test}
	repository.walk(func(source entry) {
		if source.file.IsDir() || filepath.Ext(source.path) != ".go" || strings.HasSuffix(source.path, "_test.go") {
			return
		}
		ast.Inspect(repository.parse(source.path), func(node ast.Node) bool {
			structure, valid := node.(*ast.StructType)
			if !valid {
				return true
			}
			for _, field := range structure.Fields.List {
				if repository.carries(field.Type) {
					test.Errorf("%s stores context.Context in a struct field", source.path)
				}
			}
			return true
		})
	})
}

func (repository repository) carries(expression ast.Expr) bool {
	selector, valid := expression.(*ast.SelectorExpr)
	if !valid || selector.Sel.Name != "Context" {
		return false
	}
	owner, valid := selector.X.(*ast.Ident)
	return valid && owner.Name == "context"
}
