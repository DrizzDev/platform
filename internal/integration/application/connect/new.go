package connect

import "github.com/DrizzDev/platform/internal/integration/domain/agent"

func New(options Options) (Installer, error) {
	if failure := options.validate(); failure != nil {
		return Installer{}, failure
	}
	return Installer{
		catalog:  agent.New(),
		resolver: options.Resolver,
		store:    options.Store,
		recorder: options.Recorder,
		monitor:  options.Monitor,
		logger:   options.Logger,
	}, nil
}
