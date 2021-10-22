package mCache

import (
	"errors"
	"fmt"
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

func NewCache(expiration time.Duration, cleanInterval time.Duration) *Cache {
	cache := &Cache{
		items:      make(map[string]Item),
		expiration: expiration,
		mutex:      new(sync.RWMutex),
	}
	cache.runCleanCache(cleanInterval)
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
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exist := c.items[key]
	if !exist {
		return nil, errors.New("key not exist")
	}
	if time.Now().UnixNano() > item.Expiration {
		// Если в момент запроса кеш устарел возвращаем nil
		delete(c.items, key) //и сразу удаляем его
		return nil, errors.New("expiration key")
	}
	return item.Value, nil
}

func (c *Cache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exist := c.items[key]; !exist {
		return errors.New("key not found")
	}
	delete(c.items, key)
	return nil
}

func (c *Cache) Clean() {
	//находим элементы с истекшим временем жизни и удаляем из хранилища
	if keys := c.expiredKeys(); len(keys) != 0 {
		c.clearItems(keys)
	}
}

func (c *Cache) clearItems(keys []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}

func (c *Cache) expiredKeys() (keys []string) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for i, item := range c.items {
		if time.Now().UnixNano() > item.Expiration {
			keys = append(keys, i)
		}
	}
	return
}

func (c *Cache) Size() (size int) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.items)
}

func (c *Cache) runCleanCache(interval time.Duration) {
	go func() {
		for {
			select {
			case <-time.NewTicker(interval).C:
				c.Clean()
				fmt.Println("clean, size:", c.Size()) //это только чтобы посмотреть как работает
			}
		}
	}()
}
