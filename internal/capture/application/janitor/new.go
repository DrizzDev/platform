package janitor

func New(options Options) (Janitor, error) {
	if failure := options.validate(); failure != nil {
		return Janitor{}, failure
	}
	return Janitor{
		relief:  options.low(),
		grace:   options.span(),
		ceiling: options.budget(),

		vault:     options.Vault,
		clock:     options.Clock,
		archive:   options.Archive,
		notifier:  options.Notifier,
		catalogue: options.Catalogue,
	}, nil
}
