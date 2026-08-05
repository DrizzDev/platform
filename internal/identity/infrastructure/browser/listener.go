package browser

import (
	"net/http"
	"sync/atomic"

	"github.com/DrizzDev/platform/internal/identity/application/login"
)

const misses = 3

// Interrupted reports that the loopback listener closed before a valid callback
// arrived, for example after repeated stray requests.
type Interrupted struct{}

func (Interrupted) Error() string {
	return "the sign-in callback did not arrive"
}

type capture struct {
	callbacks chan login.Callback
	faults    chan error
	misses    atomic.Int32
}

// route returns the callback handler. Only a GET to the exact callback path with
// code and state is accepted; anything else is refused, and repeated stray
// requests close the listener.
func (browser Browser) route(capture *capture) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != browser.path {
			writer.WriteHeader(http.StatusNotFound)
			if capture.misses.Add(1) >= misses {
				capture.signal(Interrupted{})
			}
			return
		}
		query := request.URL.Query()
		if request.Method != http.MethodGet || query.Get("code") == "" || query.Get("state") == "" {
			capture.reject(writer)
			capture.signal(login.Rejected{})
			return
		}
		capture.accept(writer)
		capture.deliver(login.Callback{Code: query.Get("code"), State: query.Get("state")})
	}
}

func (capture *capture) deliver(callback login.Callback) {
	select {
	case capture.callbacks <- callback:
	default:
	}
}

func (capture *capture) signal(failure error) {
	select {
	case capture.faults <- failure:
	default:
	}
}
