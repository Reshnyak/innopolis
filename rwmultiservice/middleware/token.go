package middleware

import (
	"crypto/rand"
	"fmt"
	"sync"
)

var (
	UserTokens = NewToken()
	tokenLen   = 16
)

type Token struct {
	exist map[string]struct{}
	mu    sync.RWMutex
}

func NewToken() *Token {
	return &Token{exist: make(map[string]struct{})}
}

func (t *Token) Generate() (string, error) {
	bytes := make([]byte, tokenLen)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("there were some troubles. Try again later. %s", err)
	}
	hash := string(bytes)
	t.mu.Lock()
	t.exist[hash] = struct{}{}
	t.mu.Unlock()
	return hash, nil
}
func (t *Token) IsValid(key string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.exist[key]
	return ok
}
