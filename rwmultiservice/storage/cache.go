package storage

import (
	"github.com/Reshnyak/innopolis/rwmultiservice/internal/models"
	"log"
	"sync"
)

type CacheMsg struct {
	m  map[string][]models.Message
	mu *sync.RWMutex
}

func NewCache() *CacheMsg {
	return &CacheMsg{
		m:  make(map[string][]models.Message),
		mu: new(sync.RWMutex),
	}

}
func (c *CacheMsg) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.m)

}
func (c *CacheMsg) GetKeys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.m))
	for key, _ := range c.m {
		keys = append(keys, key)
	}
	return keys

}
func (c *CacheMsg) Set(key string, value models.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = append(c.m[key], value)
}

func (c *CacheMsg) Get(key string) ([]models.Message, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.m[key]
	if ok {
		delete(c.m, key)
	}
	return val, ok

}
func (c *CacheMsg) Delete(key string) {
	log.Printf("Delete by key:%s\n", key)
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, key)
}

/// Alternative implementation
//type CacheSyncMsg struct {
//	sync.Map
//}
//
//func (c *CacheSyncMsg) Set(key string, value models.Message) {
//	messages, ok := c.Get(key)
//	if !ok {
//		messages = []models.Message{}
//	}
//	messages = append(messages, value)
//	c.Store(key, messages)
//}
//
//func (c *CacheSyncMsg) Get(key string) ([]models.Message, bool) {
//	value, ok := c.LoadAndDelete(key)
//	if ok {
//		return value.([]models.Message), ok
//	}
//	return nil, ok
//}
