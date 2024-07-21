package limiter

import (
	"sync"
)

type Limiter struct {
	cond  sync.Cond
	cache sync.Map
	limit int
}

func New(limit int) *Limiter {

	return &Limiter{
		cond:  sync.Cond{L: new(sync.Mutex)},
		limit: limit,
		cache: sync.Map{},
	}
}
func (lim *Limiter) LimitLock(key string) {
	lim.cond.L.Lock()
	for !lim.CheckAndStore(key) {
		lim.cond.Wait()
	}
	lim.cond.L.Unlock()
}
func (lim *Limiter) LimitUnLock(key string) {
	lim.cond.L.Lock()
	lim.CheckAndDelete(key)
	lim.cond.Signal()
	lim.cond.L.Unlock()
}
func (lim *Limiter) CheckAndDelete(key string) {
	if value, load := lim.cache.Load(key); load {
		if value.(int) <= 1 {
			lim.cache.Delete(key)
			return
		}
		lim.cache.Store(key, value.(int)-1)
	}
}
func (lim *Limiter) CheckAndStore(key string) bool {
	actual := 0
	if value, load := lim.cache.Load(key); load {
		if value.(int) >= lim.limit {
			return false
		}
		actual = value.(int)
	}
	lim.cache.Store(key, actual+1)
	return true
}

//type Middleware interface {
//	Handler(h http.Handler) http.Handler
//}
