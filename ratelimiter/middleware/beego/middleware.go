package beego

import (
	"net/http"
	"strings"

	"github.com/Reshnyak/innopolis/ratelimiter/limiter"
)

// Middleware is the middleware for gin.
type Middleware struct {
	Limiter *limiter.Limiter
}

// NewBeeMiddleware return a new instance of a beego middleware.
func NewBeeMiddleware(limiter *limiter.Limiter) *Middleware {

	middleware := &Middleware{
		Limiter: limiter,
	}
	return middleware
}
func (middleware *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ind := strings.LastIndex(r.RemoteAddr, ":")
		key := r.RemoteAddr[:ind]
		middleware.Limiter.LimitLock(key)
		next.ServeHTTP(w, r)
		middleware.Limiter.LimitUnLock(key)

	})
}
