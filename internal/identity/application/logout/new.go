package logout

func New(options Options) (Flow, error) {
	if failure := options.validate(); failure != nil {
		return Flow{}, failure
	}
	return Flow{
		vault:       options.Vault,
		publication: options.Publication,
		revocation:  options.Revocation,
		clock:       options.Clock,
	}, nil
}
