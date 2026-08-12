package cli

import (
	"context"

	"github.com/spf13/cobra"
)

func (command Command) root(scope context.Context) *cobra.Command {
	root := &cobra.Command{
		Use:   "drizz",
		Short: "Use Drizz capabilities from agents and developer tools",

		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.SetIn(command.options.Streams.Input)
	root.SetOut(command.options.Streams.Output)
	root.SetErr(command.options.Streams.Failure)

	root.AddCommand(command.version())
	root.AddCommand(command.mcp(scope))
	root.AddCommand(command.login(scope))
	root.AddCommand(command.logout(scope))
	command.capabilities(scope, root)

	return root
}

// capabilities adds the device-capability commands to the root, so the root builder stays small as the catalogue
// grows. The completeness gate asserts every catalogued capability is registered here.
func (command Command) capabilities(scope context.Context, root *cobra.Command) {
	root.AddCommand(command.screenshot(scope))
	root.AddCommand(command.snapshot(scope))
	root.AddCommand(command.hierarchy(scope))
	root.AddCommand(command.dimensions(scope))
	root.AddCommand(command.devices(scope))
	root.AddCommand(command.tap(scope))
	root.AddCommand(command.swipe(scope))
	root.AddCommand(command.pinch(scope))
	root.AddCommand(command.press(scope))
	root.AddCommand(command.typing(scope))
	root.AddCommand(command.clear(scope))
	root.AddCommand(command.back(scope))
	root.AddCommand(command.home(scope))
	root.AddCommand(command.locate(scope))
	root.AddCommand(command.install(scope))
	root.AddCommand(command.launch(scope))
	root.AddCommand(command.terminate(scope))
	root.AddCommand(command.wipe(scope))
	root.AddCommand(command.installed(scope))
	root.AddCommand(command.running(scope))
	root.AddCommand(command.foreground(scope))
	root.AddCommand(command.url(scope))
	root.AddCommand(command.disk(scope))
	root.AddCommand(command.images(scope))
	root.AddCommand(command.boot(scope))
	root.AddCommand(command.pause(scope))
	root.AddCommand(command.resume(scope))
}
