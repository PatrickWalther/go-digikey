package digikey

import (
	"sync"
	"time"
)

// RateLimiter tracks API usage against Digi-Key's rate limits.
// Limits: 120 requests/minute, 1000 requests/day.
type RateLimiter struct {
	mu sync.Mutex

	// Minute tracking
	minuteCount     int
	minuteResetTime time.Time

	// Day tracking
	dayCount     int
	dayResetTime time.Time

	// Limits
	minuteLimit int
	dayLimit    int
}

// NewRateLimiter creates a new rate limiter with default Digi-Key limits.
func NewRateLimiter() *RateLimiter {
	now := time.Now()
	return &RateLimiter{
		minuteLimit:     120,
		dayLimit:        1000,
		minuteResetTime: now.Add(time.Minute),
		dayResetTime:    now.Add(24 * time.Hour),
	}
}

// NewRateLimiterWithLimits creates a rate limiter with custom limits.
func NewRateLimiterWithLimits(minuteLimit, dayLimit int) *RateLimiter {
	now := time.Now()
	return &RateLimiter{
		minuteLimit:     minuteLimit,
		dayLimit:        dayLimit,
		minuteResetTime: now.Add(time.Minute),
		dayResetTime:    now.Add(24 * time.Hour),
	}
}

// Allow checks if a request is allowed and increments counters if so.
// Returns an error if the rate limit would be exceeded.
func (r *RateLimiter) Allow() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// Reset minute counter if window has passed
	if now.After(r.minuteResetTime) {
		r.minuteCount = 0
		r.minuteResetTime = now.Add(time.Minute)
	}

	// Reset day counter if window has passed
	if now.After(r.dayResetTime) {
		r.dayCount = 0
		r.dayResetTime = now.Add(24 * time.Hour)
	}

	// Check minute limit
	if r.minuteCount >= r.minuteLimit {
		return &RateLimitError{
			Limit:     r.minuteLimit,
			Remaining: 0,
			ResetAt:   r.minuteResetTime.Format(time.RFC3339),
			Type:      "minute",
		}
	}

	// Check day limit
	if r.dayCount >= r.dayLimit {
		return &RateLimitError{
			Limit:     r.dayLimit,
			Remaining: 0,
			ResetAt:   r.dayResetTime.Format(time.RFC3339),
			Type:      "day",
		}
	}

	// Increment counters
	r.minuteCount++
	r.dayCount++

	return nil
}

// Stats returns current rate limit statistics.
func (r *RateLimiter) Stats() RateLimitStats {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// Check if windows have reset
	minuteCount := r.minuteCount
	if now.After(r.minuteResetTime) {
		minuteCount = 0
	}

	dayCount := r.dayCount
	if now.After(r.dayResetTime) {
		dayCount = 0
	}

	return RateLimitStats{
		MinuteLimit:     r.minuteLimit,
		MinuteUsed:      minuteCount,
		MinuteRemaining: r.minuteLimit - minuteCount,
		MinuteResetAt:   r.minuteResetTime,
		DayLimit:        r.dayLimit,
		DayUsed:         dayCount,
		DayRemaining:    r.dayLimit - dayCount,
		DayResetAt:      r.dayResetTime,
	}
}

// RateLimitStats contains current rate limit usage information.
type RateLimitStats struct {
	MinuteLimit     int
	MinuteUsed      int
	MinuteRemaining int
	MinuteResetAt   time.Time
	DayLimit        int
	DayUsed         int
	DayRemaining    int
	DayResetAt      time.Time
}

// WaitTime returns how long to wait before the next request is allowed.
// Returns 0 if a request can be made immediately.
func (r *RateLimiter) WaitTime() time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// Get effective counts (accounting for window resets)
	minuteCount := r.minuteCount
	if now.After(r.minuteResetTime) {
		minuteCount = 0
	}

	dayCount := r.dayCount
	if now.After(r.dayResetTime) {
		dayCount = 0
	}

	// If both limits have room, no wait needed
	if minuteCount < r.minuteLimit && dayCount < r.dayLimit {
		return 0
	}

	// Calculate wait time based on which limit is exceeded
	var wait time.Duration

	if minuteCount >= r.minuteLimit && now.Before(r.minuteResetTime) {
		wait = r.minuteResetTime.Sub(now)
	}

	if dayCount >= r.dayLimit && now.Before(r.dayResetTime) {
		dayWait := r.dayResetTime.Sub(now)
		if dayWait > wait {
			wait = dayWait
		}
	}

	return wait
}

// UpdateFromResponse updates rate limiter state based on API response headers.
// Call this when receiving a 429 response with Retry-After header.
func (r *RateLimiter) UpdateFromResponse(retryAfterSeconds int) {
	if retryAfterSeconds <= 0 {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Set minute count to limit to prevent further requests
	r.minuteCount = r.minuteLimit
	r.minuteResetTime = time.Now().Add(time.Duration(retryAfterSeconds) * time.Second)
}
