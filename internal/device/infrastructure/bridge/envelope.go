package bridge

import "encoding/json"

type outbound struct {
	ID      int64  `json:"id"`
	Method  string `json:"method"`
	Version string `json:"jsonrpc"`
	Params  any    `json:"params,omitempty"`
}

type alert struct {
	Version string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  handle `json:"params"`
}

type handle struct {
	ID int64 `json:"id"`
}

type inbound struct {
	ID      int64           `json:"id"`
	Version string          `json:"jsonrpc"`
	Error   *problem        `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

type problem struct {
	Data    payload `json:"data"`
	Code    int     `json:"code"`
	Message string  `json:"message"`
}

type payload struct {
	Kind string `json:"kind"`
}

type answer struct {
	failure error
	result  json.RawMessage
}

func (inbound inbound) resolve() answer {
	if inbound.Error != nil {
		return answer{failure: inbound.Error.translate()}
	}
	return answer{result: inbound.Result}
}

// translate surfaces only the error kind; the raw sidecar message is dropped so no internal detail reaches a caller.
func (problem problem) translate() error {
	if problem.Data.Kind != "" {
		return Fault{kind: problem.Data.Kind}
	}
	return Unavailable{}
}
