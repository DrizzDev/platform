package bridge

import "context"

// spawn is the production session factory: the binary is integrity-checked before
// exec and the health handshake completes before the session is handed out.
func (channel *Channel) spawn(scope context.Context) (*session, error) {
	if failure := channel.options.verify(); failure != nil {
		return nil, failure
	}
	machine, failure := channel.options.spawn(scope)
	if failure != nil {
		return nil, failure
	}
	live := &session{
		bound:   frame,
		input:   machine.input,
		depart:  make(chan struct{}),
		stop:    machine.command.Cancel,
		outbox:  make(chan []byte, inflight),
		pending: make(map[int64]chan answer),
		reap:    func() { _ = machine.command.Wait() },
	}

	go live.absorb(machine.logs)
	go live.pump()
	go live.drain(machine.output)

	if failure := live.greet(scope); failure != nil {
		live.expire()
		return nil, failure
	}

	return live, nil
}
