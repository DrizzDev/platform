package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
)

// glide is the default swipe duration in milliseconds when the command line does not carry one.
const glide = 300

func (command Command) swipe(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Swipe)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial> <from-x> <from-y> <to-x> <to-y>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(5),
		RunE: func(root *cobra.Command, arguments []string) error {
			whole, failure := command.whole(arguments[1:])
			if failure != nil {
				return denied{message: "Coordinates must be whole numbers."}
			}
			if _, failure := command.options.Perform.Swipe(scope, operator.Drag{
				Serial:       arguments[0],
				From:         operator.Spot{X: whole[0], Y: whole[1]},
				To:           operator.Spot{X: whole[2], Y: whole[3]},
				Milliseconds: glide,
			}); failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure = fmt.Fprintln(root.OutOrStdout(), "Swiped.")
			return failure
		},
	}
}

func (command Command) pinch(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Pinch)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial> <center-x> <center-y> <start-radius> <end-radius>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(5),
		RunE: func(root *cobra.Command, arguments []string) error {
			whole, failure := command.whole(arguments[1:])
			if failure != nil {
				return denied{message: "Coordinates and radii must be whole numbers."}
			}
			if _, failure := command.options.Perform.Pinch(scope, operator.Squeeze{
				Serial: arguments[0],
				Centre: operator.Spot{X: whole[0], Y: whole[1]},
				From:   whole[2],
				To:     whole[3],
			}); failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure = fmt.Fprintln(root.OutOrStdout(), "Pinched.")
			return failure
		},
	}
}

func (command Command) press(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Press)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial> <button>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(2),
		RunE: func(root *cobra.Command, arguments []string) error {
			if _, failure := command.options.Perform.Press(scope, operator.Key{Serial: arguments[0], Button: arguments[1]}); failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure := fmt.Fprintln(root.OutOrStdout(), "Pressed.")
			return failure
		},
	}
}

func (command Command) typing(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Type)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial> <text>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(2),
		RunE: func(root *cobra.Command, arguments []string) error {
			if _, failure := command.options.Perform.Type(scope, operator.Entry{Serial: arguments[0], Text: arguments[1]}); failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure := fmt.Fprintln(root.OutOrStdout(), "Typed.")
			return failure
		},
	}
}

func (command Command) locate(scope context.Context) *cobra.Command {
	entry, _ := catalog.New().Lookup(catalog.Locate)
	locate := &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial> <latitude> <longitude>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(3),
		RunE: func(root *cobra.Command, arguments []string) error {
			lat, failure := strconv.ParseFloat(arguments[1], 64)
			if failure != nil {
				return denied{message: "The latitude must be a number."}
			}
			lon, failure := strconv.ParseFloat(arguments[2], 64)
			if failure != nil {
				return denied{message: "The longitude must be a number."}
			}
			if _, failure := command.options.Perform.Locate(scope, operator.Fix{Serial: arguments[0], Lat: lat, Lon: lon}); failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure = fmt.Fprintln(root.OutOrStdout(), "Location set.")
			return failure
		},
	}
	// A latitude or longitude may be negative; stop flag parsing after the positional arguments so a value like
	// -122.4 is not mistaken for a flag.
	locate.Flags().SetInterspersed(false)
	return locate
}

func (command Command) clear(scope context.Context) *cobra.Command {
	return command.single(scope, gesture{name: catalog.Clear, done: "Cleared.", perform: command.options.Perform.Clear})
}

func (command Command) back(scope context.Context) *cobra.Command {
	return command.single(scope, gesture{name: catalog.Back, done: "Went back.", perform: command.options.Perform.Back})
}

func (command Command) home(scope context.Context) *cobra.Command {
	return command.single(scope, gesture{name: catalog.Home, done: "Went home.", perform: command.options.Perform.Home})
}

// gesture describes a command whose only input is the target device: its catalogue name, its confirmation, and the
// capability that performs it.
type gesture struct {
	perform func(context.Context, operator.Target) (operator.Ack, error)
	name    string
	done    string
}

// single builds a device-only command from a gesture spec, so the three serial-only commands share one shape without
// dispatching by name.
func (command Command) single(scope context.Context, spec gesture) *cobra.Command {
	entry, _ := catalog.New().Lookup(spec.name)
	return &cobra.Command{
		Use:   command.slug(entry.Title()) + " <serial>",
		Short: entry.Summary(),
		Args:  cobra.ExactArgs(1),
		RunE: func(root *cobra.Command, arguments []string) error {
			if _, failure := spec.perform(scope, operator.Target{Serial: arguments[0]}); failure != nil {
				return denied{message: command.explain(failure)}
			}
			_, failure := fmt.Fprintln(root.OutOrStdout(), spec.done)
			return failure
		},
	}
}

// whole parses each argument as a whole number.
func (Command) whole(values []string) ([]int, error) {
	numbers := make([]int, 0, len(values))
	for _, value := range values {
		number, failure := strconv.Atoi(value)
		if failure != nil {
			return nil, failure
		}
		numbers = append(numbers, number)
	}
	return numbers, nil
}
