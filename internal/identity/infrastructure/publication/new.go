package publication

func New(options Options) (Publisher, error) {
	if failure := options.validate(); failure != nil {
		return Publisher{}, failure
	}
	return Publisher{
		vault:   options.Vault,
		ledger:  options.Ledger,
		random:  options.Random,
		session: options.Session,
	}, nil
}
