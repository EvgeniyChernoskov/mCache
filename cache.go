package mCache

import (
	"errors"
	"sync"
	"time"
)

var ErrItemNotFound = errors.New("cache: item not found")

type (
	Cache struct {
		cache map[interface{}]*item
		sync.RWMutex
	}

	item struct {
		Value      interface{}
		Expiration int64
	}
)

func New(cleanInterval time.Duration) *Cache {
	с := &Cache{cache: make(map[interface{}]*item)}
	go с.runCleanCache(cleanInterval)
	return с
}

func (c *Cache) Set(key interface{}, value interface{}, ttl time.Duration) {
	c.Lock()
	c.cache[key] = &item{
		Value:      value,
		Expiration: time.Now().Add(ttl).UnixNano(),
	}
	c.Unlock()
}

func (c *Cache) Get(key string) (interface{}, error) {

	c.RLock()
	item, exist := c.cache[key]
	c.RUnlock()

	if !exist {
		return nil, ErrItemNotFound
	}
	return item.Value, nil
}

func (c *Cache) runCleanCache(interval time.Duration) {
	for {
		select {
		case <-time.NewTicker(interval).C:
			c.clean()
		}
	}
}

func (c *Cache) clean() {
	c.Lock()
	for key, item := range c.cache {
		if time.Now().UnixNano() > item.Expiration {
			delete(c.cache, key)
		}
	}
	c.Unlock()
}
