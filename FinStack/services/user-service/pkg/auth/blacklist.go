package auth

import (
	"sync"
	"time"
)

type tokenBlacklist struct {
	blacklisted map[string]int64
	mu          sync.RWMutex
}

var blacklist = &tokenBlacklist{
	blacklisted: make(map[string]int64),
}

// AddToken adds a token to the blacklist with its expiration timestamp.
func AddToken(token string, exp int64) {
	blacklist.mu.Lock()
	defer blacklist.mu.Unlock()
	blacklist.blacklisted[token] = exp
}

// IsTokenBlacklisted checks if a token is blacklisted or expired.
func IsTokenBlacklisted(token string) bool {
	blacklist.mu.RLock()
	exp, exists := blacklist.blacklisted[token]
	blacklist.mu.RUnlock()
	if !exists {
		return false
	}
	return time.Now().Unix() >= exp
}
