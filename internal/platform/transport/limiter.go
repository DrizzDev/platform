package transport

import (
	"io"
	"net/http"
)

// limiter caps every response body so a malicious or misconfigured peer cannot
// exhaust memory (SEC-022).
type limiter struct {
	inner   http.RoundTripper
	ceiling int64
}

func (limiter limiter) RoundTrip(request *http.Request) (*http.Response, error) {
	response, failure := limiter.inner.RoundTrip(request)
	if failure != nil {
		return nil, failure
	}
	response.Body = capped{reader: io.LimitReader(response.Body, limiter.ceiling), body: response.Body}
	return response, nil
}

type capped struct {
	reader io.Reader
	body   io.ReadCloser
}

func (capped capped) Read(buffer []byte) (int, error) {
	return capped.reader.Read(buffer)
}

func (capped capped) Close() error {
	return capped.body.Close()
}
