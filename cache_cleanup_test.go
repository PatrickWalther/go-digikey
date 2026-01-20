package digikey

import (
	"testing"
	"time"
)

// TestMemoryCacheCleanup tests that expired entries are cleaned up
func TestMemoryCacheCleanup(t *testing.T) {
	cache := NewMemoryCache(1 * time.Second)

	// Set an entry with very short TTL
	cache.Set("short-lived", []byte("data"), 100*time.Millisecond)

	// Verify it exists
	if _, ok := cache.Get("short-lived"); !ok {
		t.Error("entry should exist immediately after set")
	}

	// Wait for it to expire
	time.Sleep(150 * time.Millisecond)

	// Entry should be gone (expired)
	if _, ok := cache.Get("short-lived"); ok {
		t.Error("expired entry should be removed")
	}
}

// TestMemoryCacheLongLivedEntry tests long-lived entries are not expired
func TestMemoryCacheLongLivedEntry(t *testing.T) {
	cache := NewMemoryCache(10 * time.Second)

	cache.Set("long-lived", []byte("data"), 5*time.Second)

	// Should still exist after a short time
	time.Sleep(100 * time.Millisecond)
	if _, ok := cache.Get("long-lived"); !ok {
		t.Error("non-expired entry should exist")
	}
}

// TestMemoryCacheCleanupWithMultipleEntries tests cleanup with multiple entries
func TestMemoryCacheCleanupWithMultipleEntries(t *testing.T) {
	cache := NewMemoryCache(1 * time.Second)

	// Set entries with different TTLs
	cache.Set("expire-early", []byte("data1"), 50*time.Millisecond)
	cache.Set("expire-late", []byte("data2"), 500*time.Millisecond)
	cache.Set("long-term", []byte("data3"), 5*time.Second)

	// Wait for early entry to expire
	time.Sleep(100 * time.Millisecond)

	// Early entry should be gone, others should exist
	if _, ok := cache.Get("expire-early"); ok {
		t.Error("early-expiring entry should be removed")
	}
	if _, ok := cache.Get("expire-late"); !ok {
		t.Error("late-expiring entry should still exist")
	}
	if _, ok := cache.Get("long-term"); !ok {
		t.Error("long-term entry should still exist")
	}
}

// TestCacheKeyForDetails tests cache key generation for product details
func TestCacheKeyForDetails(t *testing.T) {
	locale := Locale{Site: "US", Language: "en", Currency: "USD"}
	key1 := cacheKeyForDetails(locale, "PROD-123")
	key2 := cacheKeyForDetails(locale, "PROD-123")

	// Same inputs should produce same key
	if key1 != key2 {
		t.Error("cache keys should be consistent")
	}

	// Different product numbers should produce different keys
	key3 := cacheKeyForDetails(locale, "PROD-456")
	if key1 == key3 {
		t.Error("different product numbers should produce different keys")
	}

	// Different locales should produce different keys
	locale2 := Locale{Site: "CA", Language: "fr", Currency: "CAD"}
	key4 := cacheKeyForDetails(locale2, "PROD-123")
	if key1 == key4 {
		t.Error("different locales should produce different keys")
	}
}

// TestCacheKeyFormat tests that cache keys are properly formatted
func TestCacheKeyFormat(t *testing.T) {
	locale := Locale{Site: "US", Language: "en", Currency: "USD"}
	key := cacheKeyForDetails(locale, "TEST-123")

	// Key should not be empty
	if key == "" {
		t.Error("cache key should not be empty")
	}

	// Key should contain the product number
	if !stringContains(key, "TEST-123") {
		t.Errorf("cache key should contain product number, got %s", key)
	}
}

// TestMemoryCacheClearRemovesAllEntries tests that Clear removes all entries
func TestMemoryCacheClearRemovesAllEntries(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)

	cache.Set("key1", []byte("val1"), 1*time.Minute)
	cache.Set("key2", []byte("val2"), 1*time.Minute)
	cache.Set("key3", []byte("val3"), 1*time.Minute)

	if cache.Size() != 3 {
		t.Errorf("expected 3 entries, got %d", cache.Size())
	}

	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("expected 0 entries after clear, got %d", cache.Size())
	}

	// All entries should be gone
	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 should be removed after clear")
	}
}

// TestMemoryCacheDeleteRemovesEntry tests that Delete removes a specific entry
func TestMemoryCacheDeleteRemovesEntry(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)

	cache.Set("delete-me", []byte("data"), 1*time.Minute)
	if _, ok := cache.Get("delete-me"); !ok {
		t.Error("entry should exist before delete")
	}

	cache.Delete("delete-me")

	if _, ok := cache.Get("delete-me"); ok {
		t.Error("entry should be removed after delete")
	}
}

// TestMemoryCacheSizeAccuracy tests that Size returns correct count
func TestMemoryCacheSizeAccuracy(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)

	if cache.Size() != 0 {
		t.Errorf("new cache should be empty, got size %d", cache.Size())
	}

	for i := 0; i < 5; i++ {
		key := "key" + string(rune('0'+byte(i)))
		cache.Set(key, []byte("val"), 1*time.Minute)
	}

	if cache.Size() != 5 {
		t.Errorf("expected size 5, got %d", cache.Size())
	}

	// Delete one
	cache.Delete("key0")
	if cache.Size() != 4 {
		t.Errorf("expected size 4 after delete, got %d", cache.Size())
	}
}

// Helper function - checks if string contains substring
func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
