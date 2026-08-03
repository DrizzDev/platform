package architecture_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestSize(test *testing.T) {
	test.Parallel()

	repository := repository{root: filepath.Join("..", ".."), test: test}
	repository.walk(func(source entry) {
		if source.file.IsDir() || filepath.Ext(source.path) != ".go" {
			return
		}
		content, failure := os.ReadFile(source.path)
		if failure != nil {
			test.Error(failure)
			return
		}
		lines := bytes.Count(content, []byte{'\n'})
		if len(content) != 0 && content[len(content)-1] != '\n' {
			lines++
		}
		if lines > 500 {
			test.Errorf("%s has %d lines; maximum is 500", source.path, lines)
		}
	})
}
