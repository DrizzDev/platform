package reconcile

func New(options Options) (Reconciler, error) {
	if failure := options.validate(); failure != nil {
		return Reconciler{}, failure
	}
	return Reconciler{queue: options.Queue, vault: options.Vault, clock: options.Clock}, nil
}
