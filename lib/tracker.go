package lib

import (
	"net/netip"
	"sync"
	"time"
)

type RequestRecord struct {
	Timestamp time.Time
	Duration  time.Duration
}

type IPResponseTimeTracker struct {
	mu      sync.Mutex
	history map[netip.Addr][]RequestRecord
}

func NewIPResponseTimeTracker() *IPResponseTimeTracker {
	return &IPResponseTimeTracker{
		history: make(map[netip.Addr][]RequestRecord),
	}
}

// Record saves a request duration and prunes records older than the window
func (t *IPResponseTimeTracker) Record(ip netip.Addr, duration time.Duration, window time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-window)

	// Clean old entries & append the new one
	var valid []RequestRecord
	for _, r := range t.history[ip] {
		if r.Timestamp.After(cutoff) {
			valid = append(valid, r)
		}
	}
	t.history[ip] = append(valid, RequestRecord{Timestamp: now, Duration: duration})
}

// Sum calculates cumulative response time within a window for an IP address
func (t *IPResponseTimeTracker) Sum(ip netip.Addr, window time.Duration) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := time.Now().Add(-window)
	var total time.Duration
	for _, r := range t.history[ip] {
		if r.Timestamp.After(cutoff) {
			total += r.Duration
		}
	}
	return total
}
