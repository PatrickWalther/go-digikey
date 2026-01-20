package digikey

import (
	"testing"
	"time"
)

// TestNewRateLimiter tests rate limiter creation with defaults.
func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter()

	if rl == nil {
		t.Fatal("expected non-nil rate limiter")
	}

	stats := rl.Stats()
	if stats.MinuteLimit != 120 {
		t.Errorf("expected minute limit 120, got %d", stats.MinuteLimit)
	}
	if stats.DayLimit != 1000 {
		t.Errorf("expected day limit 1000, got %d", stats.DayLimit)
	}
}

// TestNewRateLimiterWithLimits tests rate limiter with custom limits.
func TestNewRateLimiterWithLimits(t *testing.T) {
	rl := NewRateLimiterWithLimits(50, 500)

	stats := rl.Stats()
	if stats.MinuteLimit != 50 {
		t.Errorf("expected minute limit 50, got %d", stats.MinuteLimit)
	}
	if stats.DayLimit != 500 {
		t.Errorf("expected day limit 500, got %d", stats.DayLimit)
	}
}

// TestRateLimiterAllow tests that requests are allowed within limits.
func TestRateLimiterAllow(t *testing.T) {
	rl := NewRateLimiterWithLimits(5, 100)

	for i := 0; i < 5; i++ {
		err := rl.Allow()
		if err != nil {
			t.Errorf("request %d should be allowed, got error: %v", i+1, err)
		}
	}

	// 6th request should fail minute limit
	err := rl.Allow()
	if err == nil {
		t.Fatal("expected rate limit error on 6th request")
	}

	if _, ok := err.(*RateLimitError); !ok {
		t.Fatalf("expected RateLimitError, got %T", err)
	}
}

// TestRateLimiterStats tests Stats method.
func TestRateLimiterStats(t *testing.T) {
	rl := NewRateLimiterWithLimits(10, 100)

	// Make some requests
	for i := 0; i < 3; i++ {
		_ = rl.Allow()
	}

	stats := rl.Stats()

	if stats.MinuteUsed != 3 {
		t.Errorf("expected minute used 3, got %d", stats.MinuteUsed)
	}
	if stats.MinuteRemaining != 7 {
		t.Errorf("expected minute remaining 7, got %d", stats.MinuteRemaining)
	}
	if stats.DayUsed != 3 {
		t.Errorf("expected day used 3, got %d", stats.DayUsed)
	}
	if stats.DayRemaining != 97 {
		t.Errorf("expected day remaining 97, got %d", stats.DayRemaining)
	}
}

// TestRateLimiterWaitTime tests WaitTime method when limits allow requests.
func TestRateLimiterWaitTime(t *testing.T) {
	rl := NewRateLimiterWithLimits(10, 100)

	waitTime := rl.WaitTime()
	if waitTime != 0 {
		t.Errorf("expected wait time 0 when under limit, got %v", waitTime)
	}
}

// TestRateLimiterWaitTimeExceeded tests WaitTime when limit is exceeded.
func TestRateLimiterWaitTimeExceeded(t *testing.T) {
	rl := NewRateLimiterWithLimits(2, 100)

	// Use up the limit
	_ = rl.Allow()
	_ = rl.Allow()

	// Next request should fail
	err := rl.Allow()
	if err == nil {
		t.Fatal("expected rate limit error")
	}

	// WaitTime should be > 0
	waitTime := rl.WaitTime()
	if waitTime <= 0 {
		t.Errorf("expected positive wait time when limit exceeded, got %v", waitTime)
	}

	if waitTime > 2*time.Minute {
		t.Errorf("wait time seems too long: %v", waitTime)
	}
}

// TestRateLimiterUpdateFromResponse tests updating rate limiter from API response.
func TestRateLimiterUpdateFromResponse(t *testing.T) {
	rl := NewRateLimiterWithLimits(10, 100)

	// Use some requests
	for i := 0; i < 5; i++ {
		_ = rl.Allow()
	}

	// API returns retry-after of 60 seconds
	rl.UpdateFromResponse(60)

	// Should now be rate limited
	err := rl.Allow()
	if err == nil {
		t.Fatal("expected rate limit error after UpdateFromResponse")
	}

	stats := rl.Stats()
	if stats.MinuteRemaining != 0 {
		t.Errorf("expected minute remaining 0 after rate limit, got %d", stats.MinuteRemaining)
	}
}

// TestRateLimiterUpdateFromResponseZero tests UpdateFromResponse with zero.
func TestRateLimiterUpdateFromResponseZero(t *testing.T) {
	rl := NewRateLimiter()

	before := rl.Stats()

	rl.UpdateFromResponse(0)

	after := rl.Stats()

	// Should be unchanged
	if before.MinuteUsed != after.MinuteUsed {
		t.Error("stats should not change with zero retry-after")
	}
}

// TestRateLimiterUpdateFromResponseNegative tests UpdateFromResponse with negative value.
func TestRateLimiterUpdateFromResponseNegative(t *testing.T) {
	rl := NewRateLimiter()

	before := rl.Stats()

	rl.UpdateFromResponse(-10)

	after := rl.Stats()

	// Should be unchanged
	if before.MinuteUsed != after.MinuteUsed {
		t.Error("stats should not change with negative retry-after")
	}
}

// TestRateLimiterMinuteWindowReset tests that minute window resets.
func TestRateLimiterMinuteWindowReset(t *testing.T) {
	t.Skip("Skipping: requires waiting for minute window to reset (~60s)")

	rl := NewRateLimiterWithLimits(2, 1000)

	// Use up minute limit
	_ = rl.Allow()
	_ = rl.Allow()

	err := rl.Allow()
	if err == nil {
		t.Fatal("expected rate limit error when minute limit exceeded")
	}

	// Wait for reset (in real test, would use fast-forward clock)
	time.Sleep(61 * time.Second)

	// Should now be allowed
	err = rl.Allow()
	if err != nil {
		t.Errorf("expected request to be allowed after minute reset, got %v", err)
	}
}

// TestRateLimitErrorMessage tests RateLimitError message.
func TestRateLimitErrorMessage(t *testing.T) {
	rle := &RateLimitError{
		Limit:     120,
		Remaining: 0,
		ResetAt:   time.Now().Format(time.RFC3339),
		Type:      "minute",
	}

	msg := rle.Error()
	if msg == "" {
		t.Fatal("expected non-empty error message")
	}

	if !contains(msg, "rate limit") {
		t.Errorf("error should mention rate limit: %s", msg)
	}
	if !contains(msg, "minute") {
		t.Errorf("error should mention minute: %s", msg)
	}
	if !contains(msg, "120") {
		t.Errorf("error should mention limit: %s", msg)
	}
}

// TestRateLimitErrorDay tests RateLimitError for day limit.
func TestRateLimitErrorDay(t *testing.T) {
	rle := &RateLimitError{
		Limit:     1000,
		Remaining: 0,
		ResetAt:   time.Now().Format(time.RFC3339),
		Type:      "day",
	}

	msg := rle.Error()
	if !contains(msg, "day") {
		t.Errorf("error should mention day: %s", msg)
	}
	if !contains(msg, "1000") {
		t.Errorf("error should mention limit: %s", msg)
	}
}
