package sqlite

import (
	"errors"
	"fmt"
	"os"
	"time"

	driver "modernc.org/sqlite"
)

// malformed and foreign are the SQLite result codes for an untrusted database: a corrupt image and a file that is not
// a database. The code is masked to its primary byte to catch extended variants.
const (
	malformed = 11
	foreign   = 26
)

func (Options) damaged(failure error) bool {
	var fault *driver.Error
	if !errors.As(failure, &fault) {
		return false
	}
	primary := fault.Code() & 0xFF
	return primary == malformed || primary == foreign
}

// quarantine renames a corrupt database and its journal aside so a fresh one can take their place.
// The damaged files are preserved for later inspection, never deleted, so nothing is silently discarded.
func (options Options) quarantine() error {
	stamp := time.Now().UnixNano()
	for _, suffix := range []string{"", "-wal", "-shm"} {
		source := options.Path + suffix
		if _, failure := os.Stat(source); errors.Is(failure, os.ErrNotExist) {
			continue
		}
		aside := fmt.Sprintf("%s.corrupt-%d%s", options.Path, stamp, suffix)
		if failure := os.Rename(source, aside); failure != nil {
			return failure
		}
	}
	return nil
}
