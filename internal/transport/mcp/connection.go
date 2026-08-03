package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/jsonrpc"
	protocol "github.com/modelcontextprotocol/go-sdk/mcp"
)

// limit bounds a single MCP frame. Frame size and validity are enforced here;
// the remaining SEC-022 limits are deferred per documents/exceptions.md EXC-001.
const limit = 1 << 20

type transport struct {
	input  io.ReadCloser
	output io.Writer
}

func (transport transport) Connect(context.Context) (protocol.Connection, error) {
	connection := &connection{
		writer:   transport.output,
		source:   transport.input,
		incoming: make(chan delivery),
		closed:   make(chan struct{}),
		done:     make(chan struct{}),
	}
	go connection.consume(bufio.NewReaderSize(transport.input, limit+1))
	return connection, nil
}

type delivery struct {
	message jsonrpc.Message
	failure error
}

type connection struct {
	writer   io.Writer
	source   io.ReadCloser
	incoming chan delivery
	closed   chan struct{}
	done     chan struct{}
	outcome  error
	shutdown sync.Once
	guard    sync.Mutex
}

func (connection *connection) consume(reader *bufio.Reader) {
	defer close(connection.done)
	for {
		message, failure := connection.frame(reader)
		select {
		case connection.incoming <- delivery{message: message, failure: failure}:
		case <-connection.closed:
			return
		}
		if failure != nil {
			return
		}
	}
}

func (connection *connection) frame(reader *bufio.Reader) (jsonrpc.Message, error) {
	for {
		line, failure := reader.ReadSlice('\n')
		if errors.Is(failure, bufio.ErrBufferFull) {
			return nil, excess{}
		}
		content := bytes.TrimRight(line, "\r\n")
		switch {
		case errors.Is(failure, io.EOF) && len(content) == 0:
			return nil, io.EOF
		case errors.Is(failure, io.EOF), len(content) > limit:
			return nil, malformed{}
		case failure != nil:
			return nil, failure
		case len(content) == 0:
			continue
		}
		if !json.Valid(content) {
			return nil, malformed{}
		}
		message, decoded := jsonrpc.DecodeMessage(content)
		if decoded != nil {
			return nil, malformed{}
		}
		return message, nil
	}
}

func (connection *connection) Read(scope context.Context) (jsonrpc.Message, error) {
	select {
	case <-scope.Done():
		return nil, scope.Err()
	case <-connection.closed:
		return nil, io.EOF
	case delivery := <-connection.incoming:
		return delivery.message, delivery.failure
	}
}

func (connection *connection) Write(scope context.Context, message jsonrpc.Message) error {
	if scope.Err() != nil {
		return scope.Err()
	}
	encoded, failure := jsonrpc.EncodeMessage(message)
	if failure != nil {
		return failure
	}
	connection.guard.Lock()
	defer connection.guard.Unlock()
	_, failure = connection.writer.Write(append(encoded, '\n'))
	return failure
}

func (connection *connection) Close() error {
	connection.shutdown.Do(func() {
		close(connection.closed)
		connection.outcome = connection.source.Close()
	})
	return connection.outcome
}

func (connection *connection) SessionID() string {
	return ""
}
