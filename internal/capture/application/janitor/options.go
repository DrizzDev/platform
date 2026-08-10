package janitor

import (
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/catalogue"
)

// Options configures a Janitor. Ceiling is the artifact disk budget; Relief is the low-water mark reclaim frees down
// to before lifting the write gate, giving the budget hysteresis; Grace spares an object young enough to still be
// racing the entry that will reference it. Each defaults in one place when left zero.
type Options struct {
	Ceiling int64
	Relief  int64
	Vault   vault
	Archive archive

	Clock     clock
	Notifier  notifier
	Grace     time.Duration
	Catalogue catalogue.Catalogue
}

func (options Options) budget() int64 {
	if options.Ceiling == 0 {
		return ceiling
	}
	return options.Ceiling
}

func (options Options) low() int64 {
	if options.Relief == 0 {
		return relief
	}
	return options.Relief
}

func (options Options) span() time.Duration {
	if options.Grace == 0 {
		return grace
	}
	return options.Grace
}
