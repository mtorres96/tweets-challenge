package http

import (
	"os"
	"strconv"
	"sync"
	"time"
)

type RateLimiter struct {
	enabled bool
	window  time.Duration
	max     int
	mu      sync.Mutex
	buckets map[string]*bucket
}

type bucket struct {
	count int
	reset time.Time
}

func NewRateLimiterFromEnv() *RateLimiter {
	enabled := true
	if v := os.Getenv("RATE_LIMIT_ENABLED"); v != "" {
		if v == "0" || v == "false" || v == "False" {
			enabled = false
		}
	}
	win := 60
	if v := os.Getenv("RATE_LIMIT_WINDOW_SEC"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			win = n
		}
	}
	max := 20
	if v := os.Getenv("RATE_LIMIT_MAX_TWEETS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			max = n
		}
	}
	return &RateLimiter{
		enabled: enabled,
		window:  time.Duration(win) * time.Second,
		max:     max,
		buckets: make(map[string]*bucket),
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	if rl == nil || !rl.enabled {
		return true
	}
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, ok := rl.buckets[key]
	if !ok || now.After(b.reset) {
		rl.buckets[key] = &bucket{count: 1, reset: now.Add(rl.window)}
		return true
	}
	if b.count < rl.max {
		b.count++
		return true
	}
	return false
}
