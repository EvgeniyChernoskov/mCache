package mCache

import (
	"errors"
	"sync"
	"time"
)

type Item struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	items      map[string]Item
	expiration time.Duration
	mutex      *sync.RWMutex
}

func New(expiration time.Duration, cleanInterval time.Duration) *Cache {
	cache := &Cache{
		items:      make(map[string]Item),
		expiration: expiration,
		mutex:      new(sync.RWMutex),
	}
	go cache.runCleanCache(cleanInterval)
	return cache
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if ttl == 0 {
		ttl = c.expiration
	}
	expiration := time.Now().Add(ttl).UnixNano()
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, exist := c.items[key]
	if !exist {
		return nil, errors.New("key not exist")
	}
	return item.Value, nil
}

func (c *Cache) Clean() {
	//находим элементы с истекшим временем жизни и удаляем из хранилища
	for key, item := range c.items {
		if time.Now().UnixNano() > item.Expiration {
			delete(c.items, key)
		}
	}
}

func (c *Cache) runCleanCache(interval time.Duration) {
	for {
		select {
		case <-time.NewTicker(interval).C:
			c.Clean()
		}
	}
}
