package ai

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

const defaultCacheSize = 100

// LRUCache is a simple thread-safe LRU cache for AI results.
type LRUCache struct {
	mu      sync.Mutex
	entries map[string]any
	order   []string
	maxSize int
}

// NewLRUCache creates a new LRU cache with the given max size.
func NewLRUCache(maxSize int) *LRUCache {
	if maxSize <= 0 {
		maxSize = defaultCacheSize
	}
	return &LRUCache{
		entries: make(map[string]any, maxSize),
		order:   make([]string, 0, maxSize),
		maxSize: maxSize,
	}
}

// CollectionKey generates a cache key from userID and collection hash.
func CollectionKey(userID uuid.UUID, collection any) (string, error) {
	data, err := json.Marshal(collection)
	if err != nil {
		return "", fmt.Errorf("marshal collection: %w", err)
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%s:%x", userID, hash), nil
}

// Get retrieves a cached value by key.
func (c *LRUCache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.entries[key]
	return v, ok
}

// Set stores a value in the cache, evicting LRU entry if at capacity.
func (c *LRUCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.entries[key]; !exists {
		if len(c.order) >= c.maxSize {
			evict := c.order[0]
			c.order = c.order[1:]
			delete(c.entries, evict)
		}
		c.order = append(c.order, key)
	}
	c.entries[key] = value
}

// InvalidateUser removes all cached entries for a user.
func (c *LRUCache) InvalidateUser(userID uuid.UUID) {
	prefix := userID.String() + ":"
	c.mu.Lock()
	defer c.mu.Unlock()

	newOrder := c.order[:0]
	for _, k := range c.order {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			delete(c.entries, k)
		} else {
			newOrder = append(newOrder, k)
		}
	}
	c.order = newOrder
}
