package cli

import (
	"context"
	"slices"
)

type Command struct {
	options Options
}

func New(options Options) (Command, error) {
	if failure := options.validate(); failure != nil {
		return Command{}, failure
	}
	options.Arguments = slices.Clone(options.Arguments)
	return Command{options: options}, nil
}

func (command Command) Run(scope context.Context) error {
	root := command.root(scope)
	root.SetArgs(command.options.Arguments)
	return root.ExecuteContext(scope)
}
