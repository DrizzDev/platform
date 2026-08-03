package architecture_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
)

func TestStandards(test *testing.T) {
	test.Parallel()

	pattern := regexp.MustCompile("(?m)^- `([A-Z]+)-([0-9]{3})`:")
	files, failure := filepath.Glob(filepath.Join("..", "..", "documents", "standards", "*.md"))
	if failure != nil {
		test.Fatal(failure)
	}
	for _, path := range files {
		content, failure := os.ReadFile(path)
		if failure != nil {
			test.Fatal(failure)
		}
		prior := make(map[string]int)
		for _, match := range pattern.FindAllSubmatch(content, -1) {
			prefix := string(match[1])
			number, failure := strconv.Atoi(string(match[2]))
			if failure != nil {
				test.Fatal(failure)
			}
			expected := prior[prefix] + 1
			if number != expected {
				test.Errorf("%s uses %s-%03d where %s-%03d is required", path, prefix, number, prefix, expected)
			}
			prior[prefix] = number
		}
	}
}

func TestIdentifiers(test *testing.T) {
	test.Parallel()

	pattern := regexp.MustCompile(`\b(ses|req|corr?|usr|org|wrk|hst|cal|exe|idm)` + `_[[:alnum:]]+\b`)
	repository := repository{root: filepath.Join("..", ".."), test: test}
	repository.walk(func(source entry) {
		if source.file.IsDir() || !repository.text(source.path) {
			return
		}
		content, failure := os.ReadFile(source.path)
		if failure != nil {
			test.Fatal(failure)
		}
		if match := pattern.Find(content); match != nil {
			test.Errorf("%s contains shortened identifier %q", source.path, match)
		}
	})
}

func (repository repository) text(path string) bool {
	switch filepath.Ext(path) {
	case ".go", ".json", ".md", ".toml", ".yaml", ".yml":
		return true
	default:
		return false
	}
}
