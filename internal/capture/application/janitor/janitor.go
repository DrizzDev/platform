package janitor

import (
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/catalogue"
)

// Janitor reclaims disk the capture stores no longer need. Each pass removes synchronized records past their
// retention and deletes the artifacts they leave unreferenced. Under disk pressure it reclaims synchronized records
// even within retention — safe, since the cloud already holds them — and degrades new writes rather than deleting
// evidence the cloud has not yet acknowledged.
type Janitor struct {
	ceiling int64
	relief  int64
	vault   vault
	archive archive

	clock     clock
	notifier  notifier
	grace     time.Duration
	catalogue catalogue.Catalogue
}
