package cache

import (
	"errors"
	"sync"
	"time"
)

type entry struct {
	value      []byte
	expiration time.Time
}

type Cache struct {
	mu sync.RWMutex
	data map[string]entry
}

func New() *Cache {
	return &Cache{
		data: make(map[string]entry),
	}
}
// Put Method
func (c *Cache) Put(key string, val []byte, ttl time.Duration) {
	c.mu.Lock() // mutex lock, uses reciever c, which is of the Cache type
	defer c.mu.Unlock() // defer ensures we unlock when the function terminates

	var exp time.Time

	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	c.data[key] = entry{
		value:	append([]byte(nil), val...),
		expiration: exp,
	}
}

// Get method 

func (c *Cache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return nil, errors.New("key not found")
	}

	if !e.expiration.IsZero() && time.Now().After(e.expiration) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return nil, errors.New("key expired")
	}
	return append([]byte(nil), e.value...), nil
}

// Delete method

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)	
}

// Janitor. Automatically cleans up expired keys every couple of seconds using a goroutine

func (c *Cache) StartJanitor(interval time.Duration, stopch <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <- ticker.C:
				c.cleanup()
			case <-stopch:
				return
			}
		}
	}()

}

func (c *Cache) cleanup() {
	now := time.Now()
	c.mu.Lock()
	for k, e := range c.data {
		if !e.expiration.IsZero() && now.After(e.expiration) {
			delete(c.data, k)
		}
	}
	c.mu.Unlock()
	
}