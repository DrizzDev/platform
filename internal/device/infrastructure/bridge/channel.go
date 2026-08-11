package bridge

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Channel is a supervised, long-lived stdio JSON-RPC conversation with the sidecar; the live session is an atomic
// pointer so steady-state dispatch is lock-free (only start/restart locks).
type Channel struct {
	build func(context.Context) (*session, error)

	slots chan struct{}
	done  chan struct{}

	stalls atomic.Int64
	live   atomic.Pointer[session]

	attempt time.Time
	mutex   sync.Mutex

	strikes int
	shut    bool
	options Options
}

func (channel *Channel) Invoke(scope context.Context, request Request) (Response, error) {
	if channel.options.Timeout > 0 {
		var cancel context.CancelFunc
		scope, cancel = context.WithTimeout(scope, channel.options.Timeout)
		defer cancel()
	}
	select {
	case channel.slots <- struct{}{}:
		defer channel.release()
	case <-scope.Done():
		return Response{}, scope.Err()
	case <-channel.done:
		return Response{}, Closed{}
	}
	live, failure := channel.ensure(scope)
	if failure != nil {
		return Response{}, failure
	}
	response, failure := live.exchange(scope, request)
	channel.observe(failure)
	return response, failure
}

func (channel *Channel) release() {
	<-channel.slots
}

// observe recycles a wedged-but-alive sidecar: patience consecutive request timeouts — not a clean death — tear the
// session down so the next call restarts it.
func (channel *Channel) observe(failure error) {
	if errors.Is(failure, context.DeadlineExceeded) {
		if channel.stalls.Add(1) >= patience {
			channel.stalls.Store(0)
			if live := channel.live.Load(); live != nil {
				live.expire()
			}
		}
		return
	}
	channel.stalls.Store(0)
}

// ensure returns the live session, starting one under single-flight if needed; repeated start failures back off so
// a broken sidecar is not hammered.
func (channel *Channel) ensure(scope context.Context) (*session, error) {
	if live := channel.alive(); live != nil {
		return live, nil
	}
	channel.mutex.Lock()
	defer channel.mutex.Unlock()
	if channel.shut {
		return nil, Closed{}
	}
	if live := channel.alive(); live != nil {
		return live, nil
	}
	if time.Now().Before(channel.attempt) {
		return nil, Unavailable{}
	}
	// The helper is long-lived and shared across requests, so it is started detached from the request that happens to
	// trigger the start; otherwise that request's completion would tear the process down and the next call would pay
	// to respawn it. Teardown is owned by Close and the timeout-recycle path through the session, not by this context.
	live, failure := channel.build(context.WithoutCancel(scope))
	if failure != nil {
		channel.penalize()
		return nil, failure
	}
	channel.strikes = 0
	channel.live.Store(live)
	return live, nil
}

// alive returns the running session, clearing it on departure. Safe without the mutex.
func (channel *Channel) alive() *session {
	live := channel.live.Load()
	if live == nil {
		return nil
	}
	select {
	case <-live.depart:
		channel.live.CompareAndSwap(live, nil)
		return nil
	default:
		return live
	}
}

func (channel *Channel) penalize() {
	channel.strikes++
	delay := min(lull<<min(channel.strikes-1, 6), ceiling)
	channel.attempt = time.Now().Add(delay)
}

func (channel *Channel) Close() error {
	channel.mutex.Lock()
	if channel.shut {
		channel.mutex.Unlock()
		return nil
	}

	channel.shut = true
	close(channel.done)

	live := channel.live.Swap(nil)
	channel.mutex.Unlock()

	if live != nil {
		live.expire()
	}
	return nil
}
