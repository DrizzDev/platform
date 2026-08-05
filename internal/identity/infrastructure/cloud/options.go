package cloud

import "net/http"

type Options struct {
	Agent    *http.Client
	Provider Provider
	Clock    Clock
	Base     string
}
