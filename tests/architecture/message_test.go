package architecture_test

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var logged = map[string]bool{
	"Debug": true, "Info": true, "Warn": true, "Error": true,
	"DebugContext": true, "InfoContext": true, "WarnContext": true, "ErrorContext": true,
}

var event = regexp.MustCompile(`^[a-z][a-z0-9]*(\.[a-z0-9]+)+$`)

func TestMessages(test *testing.T) {
	test.Parallel()

	repository := repository{root: filepath.Join("..", "..", "internal"), test: test}
	repository.walk(func(source entry) {
		if source.file.IsDir() || filepath.Ext(source.path) != ".go" || strings.HasSuffix(source.path, "_test.go") {
			return
		}
		if strings.HasSuffix(filepath.ToSlash(source.path), "/observability/provider.go") {
			return
		}
		ast.Inspect(repository.parse(source.path), func(node ast.Node) bool {
			call, valid := node.(*ast.CallExpr)
			if valid {
				repository.message(record{path: source.path, call: call})
			}
			return true
		})
	})
}

type record struct {
	path string
	call *ast.CallExpr
}

func (repository repository) message(record record) {
	selector, valid := record.call.Fun.(*ast.SelectorExpr)
	if !valid || !logged[selector.Sel.Name] {
		return
	}
	position := 0
	if strings.HasSuffix(selector.Sel.Name, "Context") {
		position = 1
	}
	if len(record.call.Args) <= position {
		return
	}
	if !repository.stable(record.call.Args[position]) {
		repository.test.Errorf("%s logs a non-stable event name; use a dotted lowercase literal", record.path)
	}
}

func (repository repository) stable(argument ast.Expr) bool {
	literal, valid := argument.(*ast.BasicLit)
	if !valid || literal.Kind != token.STRING {
		return false
	}
	value, failure := strconv.Unquote(literal.Value)
	return failure == nil && event.MatchString(value)
}
