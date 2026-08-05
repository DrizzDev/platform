package architecture_test

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestIdentity(test *testing.T) {
	test.Parallel()

	repository := repository{root: filepath.Join("..", "..", "internal"), test: test}
	repository.walk(func(source entry) {
		if source.file.IsDir() || filepath.Ext(source.path) != ".go" || strings.HasSuffix(source.path, "_test.go") {
			return
		}
		repository.restrict(source.path)
	})
}

func (repository repository) restrict(path string) {
	const confined = "/internal/identity/application/grant"
	owners := []string{
		"/internal/identity/application/session",
		"/internal/identity/infrastructure/auth0",
		"/internal/identity/infrastructure/cloud",
	}
	location := filepath.ToSlash(path)
	for _, line := range repository.parse(path).Imports {
		value := strings.Trim(line.Path.Value, `"`)
		if !strings.Contains(value, confined) {
			continue
		}
		allowed := false
		for _, owner := range owners {
			if strings.Contains(location, owner) {
				allowed = true
			}
		}
		if !allowed {
			repository.test.Errorf("%s imports the confined grant credential outside its owners", path)
		}
	}
}
