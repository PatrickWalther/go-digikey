package digikey

import (
	"net/http"
	"strings"
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

	now := time.Now()
	wait := time.Duration(retryAfterSeconds) * time.Second
	resetAt := now.Add(wait)
	r.minuteCount = r.minuteLimit
	r.minuteResetTime = resetAt
}

// UpdateFromRateLimitError updates limiter state using structured API
// rate-limit details.
func (r *RateLimiter) UpdateFromRateLimitError(rateErr *RateLimitError) {
	if rateErr == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	resetAt := parseRateLimitResetAt(rateErr.ResetAt)

	switch strings.ToLower(strings.TrimSpace(rateErr.Type)) {
	case "day":
		if rateErr.Limit > 0 {
			r.dayLimit = rateErr.Limit
		}
		r.dayCount = r.dayLimit
		if !resetAt.IsZero() && resetAt.After(now) {
			r.dayResetTime = resetAt
		}
	default:
		if rateErr.Limit > 0 {
			r.minuteLimit = rateErr.Limit
		}
		r.minuteCount = r.minuteLimit
		if !resetAt.IsZero() && resetAt.After(now) {
			r.minuteResetTime = resetAt
		}
	}
}

// UpdateFromHeaders syncs rate limiter state from API response headers.
// DigiKey includes X-BurstLimit-* (minute) and X-RateLimit-* (day) headers
// on all responses, not just 429s.
func (r *RateLimiter) UpdateFromHeaders(headers http.Header) {
	if headers == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Sync burst/minute limit
	if limit := parseHeaderInt(headers, "X-BurstLimit-Limit"); limit > 0 {
		r.minuteLimit = limit
	}
	if remaining := parseHeaderInt(headers, "X-BurstLimit-Remaining"); remaining >= 0 {
		used := r.minuteLimit - remaining
		if used < 0 {
			used = 0
		}
		r.minuteCount = used
	}
	if ts := parseHeaderResetTime(headers, "minute"); !ts.IsZero() {
		r.minuteResetTime = ts
	} else if seconds := parseHeaderResetSeconds(headers, "minute"); seconds > 0 {
		r.minuteResetTime = time.Now().Add(time.Duration(seconds) * time.Second)
	}

	// Sync day limit
	if limit := parseHeaderInt(headers, "X-RateLimit-Limit"); limit > 0 {
		r.dayLimit = limit
	}
	if remaining := parseHeaderInt(headers, "X-RateLimit-Remaining"); remaining >= 0 {
		used := r.dayLimit - remaining
		if used < 0 {
			used = 0
		}
		r.dayCount = used
	}
	if ts := parseHeaderResetTime(headers, "day"); !ts.IsZero() {
		r.dayResetTime = ts
	} else if seconds := parseHeaderResetSeconds(headers, "day"); seconds > 0 {
		r.dayResetTime = time.Now().Add(time.Duration(seconds) * time.Second)
	}
}

func parseRateLimitResetAt(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}
