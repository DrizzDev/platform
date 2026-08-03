package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DrizzDev/platform/internal/host"
)

func main() {
	os.Exit(run())
}

func run() int {
	scope, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	process, failure := host.New(host.Options{
		Arguments:   os.Args[1:],
		Environment: os.Environ(),
		Streams: host.Streams{
			Input:   os.Stdin,
			Output:  os.Stdout,
			Failure: os.Stderr,
		},
	})
	if failure != nil {
		_, _ = fmt.Fprintln(os.Stderr, failure)
		return 1
	}
	if failure := process.Run(scope); failure != nil {
		return 1
	}
	return 0
}
