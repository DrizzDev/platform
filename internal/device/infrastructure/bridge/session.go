package bridge

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"sync"
	"sync/atomic"
)

// session is one live sidecar conversation, created per process instance and torn down as a unit on death (REL-003).
type session struct {
	input io.Writer

	reap    func()
	stop    func() error
	outbox  chan []byte
	depart  chan struct{}
	pending map[int64]chan answer

	mutex    sync.Mutex
	sequence atomic.Int64

	bound int
	gone  bool
}

func (session *session) pump() {
	for {
		select {
		case <-session.depart:
			return
		case line := <-session.outbox:
			if _, failure := session.input.Write(line); failure != nil {
				session.expire()
				return
			}
		}
	}
}

func (session *session) drain(output io.Reader) {
	reader := bufio.NewReaderSize(output, 64<<10)
	for {
		line, failure := session.scan(reader)
		if failure != nil {
			session.expire()
			return
		}
		var message inbound
		if json.Unmarshal(line, &message) != nil {
			continue
		}
		session.settle(message)
	}
}

// scan reads one newline-delimited frame, bounded so a hostile or corrupt frame cannot exhaust memory (SEC-022,
// REL-004).
func (session *session) scan(reader *bufio.Reader) ([]byte, error) {
	var buffer []byte
	for {
		slice, failure := reader.ReadSlice('\n')
		buffer = append(buffer, slice...)
		if len(buffer) > session.bound {
			return nil, Oversize{}
		}
		if failure == nil {
			return buffer, nil
		}
		if errors.Is(failure, bufio.ErrBufferFull) {
			continue
		}
		return nil, failure
	}
}

func (session *session) settle(message inbound) {
	session.mutex.Lock()
	waiter, live := session.pending[message.ID]
	if live {
		delete(session.pending, message.ID)
	}
	session.mutex.Unlock()
	if live {
		waiter <- message.resolve()
	}
}

// expire tears the session down exactly once: it fails every pending caller, stops and reaps the process, then
// signals departure (REL-007, REL-010).
func (session *session) expire() {
	session.mutex.Lock()

	if session.gone {
		session.mutex.Unlock()
		return
	}

	session.gone = true
	waiters := session.pending

	session.pending = nil
	session.mutex.Unlock()

	for _, waiter := range waiters {
		waiter <- answer{failure: Unavailable{}}
	}
	if session.stop != nil {
		_ = session.stop()
	}
	if session.reap != nil {
		session.reap()
	}
	close(session.depart)
}

// absorb drains stderr so the sidecar's log pipe cannot fill and block the child.
func (session *session) absorb(logs io.Reader) {
	_, _ = io.Copy(io.Discard, logs)
}
