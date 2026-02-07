package digikey

import (
	"strings"
	"testing"
	"time"
)

// TestMemoryCacheSet tests basic cache set operation.
func TestMemoryCacheSet(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	defer cache.Close()

	key := "test:key"
	value := []byte("test value")

	cache.Set(key, value, 1*time.Minute)

	if cache.Size() != 1 {
		t.Errorf("expected cache size 1, got %d", cache.Size())
	}
}

// TestMemoryCacheGet tests basic cache get operation.
func TestMemoryCacheGet(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	defer cache.Close()

	key := "test:key"
	value := []byte("test value")

	cache.Set(key, value, 1*time.Minute)

	retrieved, ok := cache.Get(key)
	if !ok {
		t.Fatal("expected to find value in cache")
	}

	if string(retrieved) != string(value) {
		t.Errorf("expected value %s, got %s", value, retrieved)
	}
}

// TestMemoryCacheGetMissing tests cache get for missing key.
func TestMemoryCacheGetMissing(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	defer cache.Close()

	_, ok := cache.Get("nonexistent")
	if ok {
		t.Fatal("expected cache miss for nonexistent key")
	}
}

// TestMemoryCacheDelete tests cache delete operation.
func TestMemoryCacheDelete(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	defer cache.Close()

	key := "test:key"
	cache.Set(key, []byte("value"), 1*time.Minute)

	if cache.Size() != 1 {
		t.Errorf("expected cache size 1 after set")
	}

	cache.Delete(key)

	if cache.Size() != 0 {
		t.Errorf("expected cache size 0 after delete")
	}

	_, ok := cache.Get(key)
	if ok {
		t.Fatal("expected cache miss after delete")
	}
}

// TestMemoryCacheTTL tests that expired entries are not returned.
func TestMemoryCacheTTL(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	defer cache.Close()

	key := "test:key"
	cache.Set(key, []byte("value"), 100*time.Millisecond)

	// Should be available immediately
	_, ok := cache.Get(key)
	if !ok {
		t.Fatal("expected value in cache immediately after set")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	_, ok = cache.Get(key)
	if ok {
		t.Fatal("expected cache miss after TTL expiration")
	}
}

// TestMemoryCacheDefaultTTL tests that default TTL is used when zero is passed.
func TestMemoryCacheDefaultTTL(t *testing.T) {
	cache := NewMemoryCache(100 * time.Millisecond)
	defer cache.Close()

	key := "test:key"
	// Pass 0 as TTL to use default
	cache.Set(key, []byte("value"), 0)

	// Should be available immediately
	_, ok := cache.Get(key)
	if !ok {
		t.Fatal("expected value in cache immediately after set")
	}

	// Wait for default TTL expiration
	time.Sleep(150 * time.Millisecond)

	_, ok = cache.Get(key)
	if ok {
		t.Fatal("expected cache miss after default TTL expiration")
	}
}

// TestMemoryCacheMultipleEntries tests cache with multiple entries.
func TestMemoryCacheMultipleEntries(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	defer cache.Close()

	entries := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
		"key3": []byte("value3"),
	}

	for key, value := range entries {
		cache.Set(key, value, 1*time.Minute)
	}

	if cache.Size() != 3 {
		t.Errorf("expected cache size 3, got %d", cache.Size())
	}

	for key, expectedValue := range entries {
		value, ok := cache.Get(key)
		if !ok {
			t.Errorf("expected to find key %s in cache", key)
			continue
		}
		if string(value) != string(expectedValue) {
			t.Errorf("expected value %s for key %s, got %s", expectedValue, key, value)
		}
	}
}

// TestMemoryCacheClear tests clearing all cache entries.
func TestMemoryCacheClear(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	defer cache.Close()

	// Add multiple entries
	for i := 0; i < 5; i++ {
		cache.Set("key"+string(rune('0'+i)), []byte("value"), 1*time.Minute)
	}

	if cache.Size() == 0 {
		t.Fatal("expected cache to have entries")
	}

	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("expected cache size 0 after clear, got %d", cache.Size())
	}
}

// TestMemoryCacheSize tests the Size method.
func TestMemoryCacheSize(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	defer cache.Close()

	if cache.Size() != 0 {
		t.Errorf("expected initial cache size 0, got %d", cache.Size())
	}

	for i := 0; i < 10; i++ {
		cache.Set("key"+string(rune('0'+i)), []byte("value"), 1*time.Minute)
		expected := i + 1
		if cache.Size() != expected {
			t.Errorf("expected cache size %d, got %d", expected, cache.Size())
		}
	}
}

// TestMemoryCacheOverwrite tests overwriting existing cache entries.
func TestMemoryCacheOverwrite(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	defer cache.Close()

	key := "test:key"

	cache.Set(key, []byte("value1"), 1*time.Minute)
	if cache.Size() != 1 {
		t.Fatal("expected cache size 1 after first set")
	}

	// Overwrite with new value
	cache.Set(key, []byte("value2"), 1*time.Minute)
	if cache.Size() != 1 {
		t.Fatal("expected cache size 1 after overwrite (should not duplicate)")
	}

	value, ok := cache.Get(key)
	if !ok {
		t.Fatal("expected to find value in cache")
	}

	if string(value) != "value2" {
		t.Errorf("expected new value value2, got %s", value)
	}
}

// TestMemoryCacheEmptyValue tests storing empty values.
func TestMemoryCacheEmptyValue(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	defer cache.Close()

	key := "test:key"
	cache.Set(key, []byte(""), 1*time.Minute)

	value, ok := cache.Get(key)
	if !ok {
		t.Fatal("expected to find empty value in cache")
	}

	if len(value) != 0 {
		t.Errorf("expected empty value, got %v", value)
	}
}

// TestCacheInterface tests that MemoryCache implements Cache interface.
func TestCacheInterface(t *testing.T) {
	var _ Cache = (*MemoryCache)(nil)
}

// TestMemoryCacheClose tests that Close stops the cleanup goroutine.
func TestMemoryCacheClose(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	cache.Set("key", []byte("value"), 1*time.Minute)

	// Close should not panic
	cache.Close()

	// Cache should still be readable after Close (just no more cleanup)
	if _, ok := cache.Get("key"); !ok {
		t.Error("expected to still read cached data after Close")
	}
}

// TestMemoryCacheCleanup tests that expired entries are cleaned up
func TestMemoryCacheCleanup(t *testing.T) {
	cache := NewMemoryCache(1 * time.Second)
	defer cache.Close()

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
	defer cache.Close()

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
	defer cache.Close()

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
	if !strings.Contains(key, "TEST-123") {
		t.Errorf("cache key should contain product number, got %s", key)
	}
}

// TestMemoryCacheClearRemovesAllEntries tests that Clear removes all entries
func TestMemoryCacheClearRemovesAllEntries(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	defer cache.Close()

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
	defer cache.Close()

	cache.Set("delete-me", []byte("data"), 1*time.Minute)
	if _, ok := cache.Get("delete-me"); !ok {
		t.Error("entry should exist before delete")
	}

	cache.Delete("delete-me")

	if _, ok := cache.Get("delete-me"); ok {
		t.Error("entry should be removed after delete")
	}
}

// TestDefaultCacheConfigLookupTTL tests that LookupTTL defaults to 15 minutes.
func TestDefaultCacheConfigLookupTTL(t *testing.T) {
	config := DefaultCacheConfig()
	if config.LookupTTL != 15*time.Minute {
		t.Errorf("expected LookupTTL 15m, got %v", config.LookupTTL)
	}
}

// TestCacheKeyForLookup tests cache key generation for lookup endpoints.
func TestCacheKeyForLookup(t *testing.T) {
	locale := Locale{Site: "US", Language: "en", Currency: "USD"}

	// Same inputs produce same key
	key1 := cacheKeyForLookup(locale, "categories", "all")
	key2 := cacheKeyForLookup(locale, "categories", "all")
	if key1 != key2 {
		t.Error("cache keys should be consistent")
	}

	// Different prefix produces different key
	key3 := cacheKeyForLookup(locale, "manufacturers", "all")
	if key1 == key3 {
		t.Error("different prefixes should produce different keys")
	}

	// Different identifier produces different key
	key4 := cacheKeyForLookup(locale, "categories", "42")
	if key1 == key4 {
		t.Error("different identifiers should produce different keys")
	}

	// Different locale produces different key
	locale2 := Locale{Site: "DE", Language: "de", Currency: "EUR"}
	key5 := cacheKeyForLookup(locale2, "categories", "all")
	if key1 == key5 {
		t.Error("different locales should produce different keys")
	}

	// Key format contains prefix and identifier
	if !strings.Contains(key1, "categories") || !strings.Contains(key1, "all") {
		t.Errorf("key should contain prefix and identifier, got %s", key1)
	}
}

// TestMemoryCacheSizeAccuracy tests that Size returns correct count
func TestMemoryCacheSizeAccuracy(t *testing.T) {
	cache := NewMemoryCache(5 * time.Minute)
	defer cache.Close()

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
