package mCache

import (
	"errors"
	"time"
)

type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

type Cache struct {
	items      map[string]Item
	expiration time.Duration
}

func NewCache(expiration time.Duration) *Cache {
	cache := Cache{
		items:      make(map[string]Item),
		expiration: expiration,
	}
	return &cache
}

func (c *Cache) Set(key string, value interface{}) {
	var expiration int64
	expiration = time.Now().Add(c.expiration).UnixNano()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	item, exist := c.items[key]
	if !exist {
		return nil, false
	}
	// Если в момент запроса кеш устарел возвращаем nil
	if time.Now().UnixNano() > item.Expiration {
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) Delete(key string) error {
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
	for _, k := range keys {
		delete(c.items, k)
	}
}

func (c *Cache) expiredKeys() (keys []string) {
	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration {
			keys = append(keys, k)
		}
	}
	return
}

func (c Cache) Size() (size int){
	return len(c.items)
}