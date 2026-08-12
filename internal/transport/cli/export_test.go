package cli

import "context"

// Names exposes the command-line subcommand names so the completeness gate can assert every catalogued capability has
// a command.
func Names(options Options) ([]string, error) {
	command, failure := New(options)
	if failure != nil {
		return nil, failure
	}
	root := command.root(context.Background())
	names := make([]string, 0, len(root.Commands()))
	for _, child := range root.Commands() {
		names = append(names, child.Name())
	}
	return names, nil
}

// Slug exposes the outward-name to command-line-name mapping so the gate checks the same mapping the commands use.
func Slug(title string) string {
	return Command{}.slug(title)
}
