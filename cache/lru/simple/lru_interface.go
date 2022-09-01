package simple

// LRUCache is the interface for simple LRU cache.
type LRUCache interface {
	// Set a value to the cache, returns true if an eviction occurred and
	// updates the "recently used"-ness of the key.
	Set(key, value interface{}) bool

	// Get Returns key's value from the cache and
	// updates the "recently used"-ness of the key. #value, isFound
	Get(key interface{}) (value interface{}, ok bool)

	// Contains if a key exists in cache without updating the recent-ness.
	Contains(key interface{}) (ok bool)

	// Peek Returns key's value without updating the "recently used"-ness of the key.
	Peek(key interface{}) (value interface{}, ok bool)

	// Remove a key from the cache.
	Remove(key interface{}) bool

	// RemoveOldest Removes the oldest entry from cache.
	RemoveOldest() (interface{}, interface{}, bool)

	// GetOldest Returns the oldest entry from the cache. #key, value, isFound
	GetOldest() (interface{}, interface{}, bool)

	// Keys Returns a slice of the keys in the cache, from oldest to newest.
	Keys() []interface{}

	// Len Returns the number of items in the cache.
	Len() int

	// Purge Clears all cache entries.
	Purge()

	// Resize Resizes cache, returning number evicted
	Resize(int) int
}
