package cli

import "context"

type Runner interface {
	Run(context.Context) error
}
