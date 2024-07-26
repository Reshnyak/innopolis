package gin

import (
	"github.com/Reshnyak/innopolis/ratelimiter/limiter"
	"github.com/gin-gonic/gin"
)

// Middleware is the middleware for gin.
type Middleware struct {
	Limiter *limiter.Limiter
}

// NewGinMiddleware return a new instance of a gin middleware.
func NewGinMiddleware(limiter *limiter.Limiter) gin.HandlerFunc {
	middleware := &Middleware{
		Limiter: limiter,
	}

	return func(ctx *gin.Context) {
		middleware.Handle(ctx)
	}
}

// Handle for gin request.
func (middleware *Middleware) Handle(context *gin.Context) {
	key := context.RemoteIP()
	middleware.Limiter.LimitLock(key)
	context.Next()
	middleware.Limiter.LimitUnLock(key)

}
