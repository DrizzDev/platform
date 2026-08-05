package session

func New(options Options) (Flow, error) {
	if failure := options.validate(); failure != nil {
		return Flow{}, failure
	}
	return Flow{
		vault:       options.Vault,
		refresh:     options.Refresh,
		publication: options.Publication,
		epoch:       options.Epoch,
		clock:       options.Clock,
	}, nil
}
