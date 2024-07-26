package stdhttp

import (
	"net"
	"net/http"

	"github.com/Reshnyak/innopolis/ratelimiter/limiter"
)

type Middleware struct {
	Limiter *limiter.Limiter
}

func NewStdHttpMiddleware(limiter *limiter.Limiter) *Middleware {
	middleware := &Middleware{
		Limiter: limiter,
	}
	return middleware
}

// Handler handles a HTTP request.
func (mid *Middleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		mid.Limiter.LimitLock(key)
		h.ServeHTTP(w, r)
		mid.Limiter.LimitUnLock(key)
	})
}
