package ttlcache

// stdlib
import (
	"errors"
	"sync"
	"time"
)

var ERR_KEY_EXISTS = errors.New("TTLCache: Key exists")
var ERR_KEY_NO_EXISTS = errors.New("TTLCache: Key doesn't exist")

func NewTTLCache(ttl time.Duration) *TTLCache {
	return &TTLCache{
		ttl:  ttl,
		keys: map[string]interface{}{},
	}
}

type TTLCache struct {
	sync.Mutex
	ttl  time.Duration
	keys map[string]interface{}
}

func (c *TTLCache) Exists(key string) bool {
	c.Lock()
	defer c.Unlock()
	_, ok := c.keys[key]
	return ok
}

func (c *TTLCache) Get(key string) (interface{}, error) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.keys[key]; !ok {
		return nil, ERR_KEY_NO_EXISTS
	}
	return c.keys[key], nil
}

func (c *TTLCache) Set(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.keys[key]; ok {
		return ERR_KEY_EXISTS
	}
	c.keys[key] = value
	time.AfterFunc(c.ttl, func() {
		c.Expire(key)
	})
	return nil
}

func (c *TTLCache) Expire(key string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.keys[key]; !ok {
		return ERR_KEY_NO_EXISTS
	}
	delete(c.keys, key)
	return nil
}
