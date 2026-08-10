package bridge

import "encoding/json"

type Request struct {
	Params any
	Method string
}

// Response is the raw result; the channel never interprets it — the Bridge translation validates and maps it to a
// domain type.
type Response struct {
	Result json.RawMessage
}

func (request Request) compose(id int64) ([]byte, error) {
	line, failure := json.Marshal(outbound{Version: version, ID: id, Method: request.Method, Params: request.Params})
	if failure != nil {
		return nil, failure
	}
	return append(line, '\n'), nil
}
