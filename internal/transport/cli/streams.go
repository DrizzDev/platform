package cli

import "io"

type Streams struct {
	Input   io.Reader
	Output  io.Writer
	Failure io.Writer
}
