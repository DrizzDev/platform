package cli

import "github.com/DrizzDev/platform/internal/application/release"

type Options struct {
	MCP       Runner
	Login     Session
	Device    Session
	Logout    Departure
	Perform   Perform
	Arguments []string
	Streams   Streams
	Release   release.Identity
}
