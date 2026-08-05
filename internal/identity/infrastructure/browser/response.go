package browser

import (
	"io"
	"net/http"
)

const (
	accepted = `<!doctype html><html lang="en"><head><meta charset="utf-8"><title>Drizz</title></head>` +
		`<body><p>You are signed in to Drizz. You can close this window.</p></body></html>`
	refused = `<!doctype html><html lang="en"><head><meta charset="utf-8"><title>Drizz</title></head>` +
		`<body><p>Drizz sign-in could not be completed.</p></body></html>`
)

func (capture *capture) accept(writer http.ResponseWriter) {
	capture.harden(writer)
	writer.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(writer, accepted)
}

func (capture *capture) reject(writer http.ResponseWriter) {
	capture.harden(writer)
	writer.WriteHeader(http.StatusBadRequest)
	_, _ = io.WriteString(writer, refused)
}

// harden applies the fixed no-store, no-referrer, script-free response policy so
// the callback page reflects no provider or query content (SEC-009).
func (*capture) harden(writer http.ResponseWriter) {
	header := writer.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Cache-Control", "no-store")
	header.Set("Content-Security-Policy", "default-src 'none'")
	header.Set("Referrer-Policy", "no-referrer")
	header.Set("X-Content-Type-Options", "nosniff")
}
