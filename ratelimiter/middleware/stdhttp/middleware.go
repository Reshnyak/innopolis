package stdhttp

import (
	"github.com/Reshnyak/innopolis/ratelimiter/limiter"
	"net/http"
	"strings"
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
		ind := strings.LastIndex(r.RemoteAddr, ":")
		key := r.RemoteAddr[:ind]
		mid.Limiter.LimitLock(key)
		h.ServeHTTP(w, r)
		mid.Limiter.LimitUnLock(key)
	})
}
