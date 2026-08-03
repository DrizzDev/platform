package mcp

import (
	"bufio"
	"context"
	"io"
	"time"

	"github.com/modelcontextprotocol/go-sdk/jsonrpc"
	protocol "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Request struct {
	Server    Server
	Transport protocol.Transport
}

func Serve(scope context.Context, request Request) error {
	return request.Server.serve(scope, request.Transport)
}

func Read(input io.Reader) (jsonrpc.Message, error) {
	connection := &connection{}
	return connection.frame(bufio.NewReaderSize(input, limit+1))
}

func Rejected(failure error) bool {
	return execution{}.rejected(failure)
}

func Terminate(input io.ReadCloser) bool {
	link, _ := transport{input: input, output: io.Discard}.Connect(context.Background())
	if failure := link.Close(); failure != nil {
		return false
	}
	select {
	case <-link.(*connection).done:
		return true
	case <-time.After(2 * time.Second):
		return false
	}
}

func Discard(input io.ReadCloser) (error, error) {
	link, _ := transport{input: input, output: io.Discard}.Connect(context.Background())
	return link.Close(), link.Close()
}
