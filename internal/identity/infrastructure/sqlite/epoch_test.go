package sqlite_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
)

func TestEpoch(test *testing.T) {
	test.Parallel()

	database := fixture{test: test}.open(filepath.Join(test.TempDir(), "identity.db"))
	value, failure := database.Epoch(context.Background())
	if failure != nil {
		test.Fatal(failure)
	}
	if value != epoch.Epoch(0) {
		test.Fatalf("initial epoch = %d", value)
	}
}
