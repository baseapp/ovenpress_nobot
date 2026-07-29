package lib

import (
	"net/netip"
	"testing"
	"time"
)

func TestIPResponseTimeTracker(t *testing.T) {
	tracker := NewIPResponseTimeTracker()
	ip := netip.MustParseAddr("192.168.1.1")

	// 1. Initially sum should be 0
	if sum := tracker.Sum(ip, time.Second); sum != 0 {
		t.Errorf("expected sum to be 0, got %v", sum)
	}

	// 2. Record some requests
	tracker.Record(ip, 100*time.Millisecond, time.Second)
	tracker.Record(ip, 200*time.Millisecond, time.Second)

	if sum := tracker.Sum(ip, time.Second); sum != 300*time.Millisecond {
		t.Errorf("expected sum to be 300ms, got %v", sum)
	}

	// 3. Prune check (simulated via shorter window)
	// We'll record a request with 500ms duration
	tracker.Record(ip, 500*time.Millisecond, time.Second)
	
	// Wait a tiny bit (not full second) or simulate by calling Record which prunes older ones.
	// Since we can't control time easily without mocking, let's just make sure Sum filters by cutoff.
	time.Sleep(10 * time.Millisecond)
	
	// Sum with a tiny window that shouldn't include anything if we wait, but since we didn't wait 1s, they are still within 1s.
	if sum := tracker.Sum(ip, 10*time.Millisecond); sum != 0 {
		t.Errorf("expected sum for a 10ms window to be 0 (since they happened more than 10ms ago), got %v", sum)
	}
}
