package http

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRateLimiter_WindowAndMax(t *testing.T) {
	rl := &RateLimiter{
		enabled: true,
		window:  50 * time.Millisecond,
		max:     2,
		buckets: make(map[string]*bucket),
	}
	key := "u1"
	if !rl.Allow(key) {
		t.Fatal("first should allow")
	}
	if !rl.Allow(key) {
		t.Fatal("second should allow")
	}
	if rl.Allow(key) {
		t.Fatal("third should block")
	}

	// Esperar reinicio de ventana
	time.Sleep(60 * time.Millisecond)
	if !rl.Allow(key) {
		t.Fatal("should allow after window reset")
	}
}
func TestRateLimiter_Concurrent(t *testing.T) {
	rl := &RateLimiter{enabled: true, window: 50 * time.Millisecond, max: 100, buckets: make(map[string]*bucket)}
	key := "u1"
	const N = 200
	var allowed int64
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			if rl.Allow(key) {
				atomic.AddInt64(&allowed, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if allowed > int64(rl.max) {
		t.Fatalf("allowed=%d exceeds max=%d", allowed, rl.max)
	}
}
