package ttlcache

// stdlib
import (
	"errors"
	"sync"
	"time"
)

// ERR_KEY_NO_EXISTS is the error if an expected key is not present
var ERR_KEY_NO_EXISTS = errors.New("TTLCache: Key doesn't exist")

// NewTTLCache creates a ttl cache with a default duration
func NewTTLCache(defaultTTL time.Duration) *TTLCache {
	return &TTLCache{
		defaultTTL: defaultTTL,
		keys:       map[string]interface{}{},
	}
}

// TTLCache stores the cache state
type TTLCache struct {
	sync.Mutex
	defaultTTL time.Duration
	keys       map[string]interface{}
}

// Exists checks if a key is in the cache
func (c *TTLCache) Exists(key string) bool {
	c.Lock()
	defer c.Unlock()
	_, ok := c.keys[key]
	return ok
}

// Get the value at a key
func (c *TTLCache) Get(key string) (interface{}, error) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.keys[key]; !ok {
		return nil, ERR_KEY_NO_EXISTS
	}
	return c.keys[key], nil
}

// Set sets a key with the default ttl
func (c *TTLCache) Set(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()
	c.keys[key] = value
	time.AfterFunc(c.defaultTTL, func() {
		c.Expire(key)
	})
	return nil
}

// SetEx sets a key with a specific ttl
func (c *TTLCache) SetEx(key string, value interface{}, expireAt time.Duration) error {
	c.Lock()
	defer c.Unlock()
	c.keys[key] = value
	time.AfterFunc(expireAt, func() {
		c.Expire(key)
	})
	return nil
}

// Expire a key immediately
func (c *TTLCache) Expire(key string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.keys[key]; !ok {
		return ERR_KEY_NO_EXISTS
	}
	delete(c.keys, key)
	return nil
}
