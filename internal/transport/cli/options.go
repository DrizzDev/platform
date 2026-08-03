package cli

import "github.com/DrizzDev/platform/internal/application/release"

type Options struct {
	MCP       Runner
	Arguments []string
	Streams   Streams
	Release   release.Identity
}
