package host

import "io"

type Streams struct {
	Input   io.ReadCloser
	Output  io.Writer
	Failure io.Writer
}
