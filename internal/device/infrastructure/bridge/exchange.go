package bridge

import (
	"context"
	"encoding/json"
)

type ticket struct {
	waiter chan answer
	id     int64
}

// exchange runs one multiplexed request/reply. On cancellation it drops the pending entry and sends $/cancel so the
// sidecar can abort; a late reply is discarded (REL-001).
func (session *session) exchange(scope context.Context, request Request) (Response, error) {
	id := session.sequence.Add(1)
	waiter := make(chan answer, 1)

	if !session.enlist(ticket{id: id, waiter: waiter}) {
		return Response{}, Unavailable{}
	}

	line, failure := request.compose(id)
	if failure != nil {
		session.discharge(id)
		return Response{}, failure
	}

	select {
	case session.outbox <- line:
	case <-scope.Done():
		session.discharge(id)
		session.cancel(id)
		return Response{}, scope.Err()
	case <-session.depart:
		session.discharge(id)
		return Response{}, Unavailable{}
	}

	select {
	case result := <-waiter:
		return Response{Result: result.result}, result.failure
	case <-scope.Done():
		session.discharge(id)
		session.cancel(id)
		return Response{}, scope.Err()
	case <-session.depart:
		return Response{}, Unavailable{}
	}
}

func (session *session) enlist(ticket ticket) bool {
	session.mutex.Lock()
	defer session.mutex.Unlock()

	if session.gone {
		return false
	}
	session.pending[ticket.id] = ticket.waiter
	return true
}

func (session *session) discharge(id int64) {
	session.mutex.Lock()
	delete(session.pending, id)
	session.mutex.Unlock()
}

func (session *session) cancel(id int64) {
	line, failure := json.Marshal(alert{Version: version, Method: revoke, Params: handle{ID: id}})
	if failure != nil {
		return
	}
	line = append(line, '\n')

	select {
	case session.outbox <- line:
	case <-session.depart:
	default:
	}
}
