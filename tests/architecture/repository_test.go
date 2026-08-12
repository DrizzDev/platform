package architecture_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"testing"
)

type repository struct {
	root string
	test *testing.T
}

type entry struct {
	path string
	file fs.DirEntry
}

func (repository repository) walk(visitor func(entry)) {
	repository.test.Helper()
	failure := filepath.WalkDir(repository.root, func(path string, file fs.DirEntry, failure error) error {
		if failure != nil {
			return failure
		}
		// Skip version control and generated build output; the gates police tracked source, not release artifacts.
		if file.IsDir() && (file.Name() == ".git" || file.Name() == "dist") {
			return filepath.SkipDir
		}
		visitor(entry{path: path, file: file})
		return nil
	})
	if failure != nil {
		repository.test.Fatal(failure)
	}
}

func (repository repository) parse(path string) *ast.File {
	repository.test.Helper()
	file, failure := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if failure != nil {
		repository.test.Fatal(failure)
	}
	return file
}
