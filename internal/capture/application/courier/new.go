package courier

func New(options Options) (Courier, error) {
	if failure := options.validate(); failure != nil {
		return Courier{}, failure
	}
	return Courier{ledger: options.Ledger, vault: options.Vault, uploader: options.Uploader}, nil
}
