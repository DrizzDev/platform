package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
)

func (command Command) install(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Install)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial> <path>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(2),
		RunE: func(root *cobra.Command, arguments []string) error {
			if _, failure := command.options.Perform.Install(scope, operator.Package{Serial: arguments[0], Path: arguments[1]}); failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure := fmt.Fprintln(root.OutOrStdout(), "Installed.")
			return failure
		},
	}
}

func (command Command) launch(scope context.Context) *cobra.Command {
	return command.application(scope, appspec{name: catalog.Launch, done: "Launched.", perform: command.options.Perform.Launch})
}

func (command Command) terminate(scope context.Context) *cobra.Command {
	return command.application(scope, appspec{name: catalog.Terminate, done: "Terminated.", perform: command.options.Perform.Terminate})
}

func (command Command) wipe(scope context.Context) *cobra.Command {
	return command.application(scope, appspec{name: catalog.Wipe, done: "Cleared.", perform: command.options.Perform.Wipe})
}

// appspec describes an application action naming an app: the catalogue name, its confirmation, and the capability.
type appspec struct {
	perform func(context.Context, operator.Application) (operator.Ack, error)
	name    string
	done    string
}

func (command Command) application(scope context.Context, spec appspec) *cobra.Command {
	entry, _ := catalog.New().Lookup(spec.name)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial> <app>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(2),
		RunE: func(root *cobra.Command, arguments []string) error {
			if _, failure := spec.perform(scope, operator.Application{Serial: arguments[0], App: arguments[1]}); failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure := fmt.Fprintln(root.OutOrStdout(), spec.done)
			return failure
		},
	}
}

func (command Command) installed(scope context.Context) *cobra.Command {
	return command.catalogue(scope, listspec{name: catalog.Installed, perform: command.options.Perform.Installed})
}

func (command Command) running(scope context.Context) *cobra.Command {
	return command.catalogue(scope, listspec{name: catalog.Running, perform: command.options.Perform.Running})
}

// listspec describes an application-listing read: its catalogue name and the capability that reads it.
type listspec struct {
	perform func(context.Context, operator.Target) (operator.Listing, error)
	name    string
}

func (command Command) catalogue(scope context.Context, spec listspec) *cobra.Command {
	entry, _ := catalog.New().Lookup(spec.name)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(1),
		RunE: func(root *cobra.Command, arguments []string) error {
			listing, failure := spec.perform(scope, operator.Target{Serial: arguments[0]})
			if failure != nil {
				return denied{message: command.explain(failure)}
			}
			for _, application := range listing.Apps {
				if _, failure := fmt.Fprintf(root.OutOrStdout(), "%s\t%s\t%s\n", application.Id, application.Name, application.Note); failure != nil {
					return failure
				}
			}
			return nil
		},
	}
}

func (command Command) foreground(scope context.Context) *cobra.Command {
	return command.reading(scope, readspec{name: catalog.Foreground, perform: command.options.Perform.Foreground})
}

func (command Command) url(scope context.Context) *cobra.Command {
	return command.reading(scope, readspec{name: catalog.Url, perform: command.options.Perform.Url})
}

// readspec describes a single-value read: its catalogue name and the capability that reads it.
type readspec struct {
	perform func(context.Context, operator.Target) (operator.Report, error)
	name    string
}

func (command Command) reading(scope context.Context, spec readspec) *cobra.Command {
	entry, _ := catalog.New().Lookup(spec.name)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(1),
		RunE: func(root *cobra.Command, arguments []string) error {
			report, failure := spec.perform(scope, operator.Target{Serial: arguments[0]})
			if failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure = fmt.Fprintln(root.OutOrStdout(), report.Text)
			return failure
		},
	}
}

func (command Command) disk(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Disk)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(1),
		RunE: func(root *cobra.Command, arguments []string) error {
			measure, failure := command.options.Perform.Disk(scope, operator.Target{Serial: arguments[0]})
			if failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure = fmt.Fprintf(root.OutOrStdout(), "%d MB\n", measure.Value)
			return failure
		},
	}
}
