package ttlcache

// stdlib
import (
	"runtime"
	"testing"
	"time"
)

func TestTTLCache(t *testing.T) {
	var ttl = 1 * time.Millisecond
	t.Parallel()
	cache := TTLCache{
		defaultTTL: ttl,
		keys:       map[string]interface{}{},
	}
	setValue := interface{}("1")
	err := cache.Set("a", setValue)
	if err != nil {
		t.Fatal("Unable to set key on cache", err)
	}
	exists := cache.Exists("a")
	if !exists {
		t.Fatal("Set key does not exist")
	}
	gotValue, err := cache.Get("a")
	if err != nil {
		t.Fatal("Unable to get key on cache", err)
	}
	if gotValue != setValue {
		t.Fatal("Value of key is not what was set")
	}
	time.Sleep(ttl)
	runtime.Gosched() // Yield to other routines
	exists = cache.Exists("a")
	if exists {
		t.Fatal("Key still exists after waiting", ttl)
	}
	gotValue, err = cache.Get("a")
	if err != ERR_KEY_NO_EXISTS {
		t.Fatal("Should have gotten ERR_KEY_NO_EXISTS error", err)
	}
}
