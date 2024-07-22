package util

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// ClientLimiter contains the rate limiters for each IP
type ClientLimiter struct {
	limiterMinute *rate.Limiter
	limiterDaily  *rate.Limiter
}

// RateLimiter struct to manage the rate limiters for all clients
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*ClientLimiter
}

var (
	minuteLimit = 2
	dailyLimit  = 10
)

// NewRateLimiter initializes a new RateLimiter
func NewRateLimiter(ml int, dl int) *RateLimiter {
	minuteLimit = ml
	dailyLimit = dl
	return &RateLimiter{
		clients: make(map[string]*ClientLimiter),
	}
}

// getClientLimiter returns the rate limiter for a given IP
func (rl *RateLimiter) getClientLimiter(ip string) *ClientLimiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.clients[ip]
	if !exists {
		limiter = &ClientLimiter{
			limiterMinute: rate.NewLimiter(rate.Every(time.Minute), minuteLimit),
			limiterDaily:  rate.NewLimiter(rate.Every(24*time.Hour), dailyLimit),
		}
		rl.clients[ip] = limiter
	}
	return limiter
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	clientLimiter := rl.getClientLimiter(ip)
	return clientLimiter.limiterMinute.Allow() && clientLimiter.limiterDaily.Allow()
}
