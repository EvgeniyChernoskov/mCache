package main

import (
	"fmt"
	"github.com/EvgeniyChernoskov/mCache"
	"time"
)

func main() {
	//создаем кеш
	cache := mCache.NewCache(time.Second * 5)
	//добавляем две записи
	cache.Set("aaa", time.Now())
	cache.Set("bbb", 156)

	//читаем записи
	fmt.Println(cache.Get("aaa"))
	fmt.Println(cache.Get("bbb"))

	//читаем несуществующую запись
	value, ok := cache.Get("ccc")
	if ok {
		fmt.Println(value, value)
	} else{
		fmt.Println("not found")
	}

	//размер кеша
	printCacheSize(cache)

	//удаляем записи
	err := cache.Delete("aaa")
	if err!= nil {
		fmt.Println("aaa","deleted")
	}

	//по второй неделе занятий можно сделать очистку кеша с периодом в горутине
	//чистка раз в секунду:
	go func(cache *mCache.Cache, d time.Duration) {
		for {
			select {
			case <-time.NewTicker(d).C:
				cache.Clean()
				printCacheSize(cache)
			}
		}
	}(cache, time.Second)

	//накидаем в кеш чего нибудь
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("item№%d", i)
		value := fmt.Sprintf("value%d", i)
		cache.Set(key,value)
		fmt.Println("add:",key,value)
		time.Sleep(time.Second)
	}
	printCacheSize(cache)


	//ну и чтобы было видно что происходит
	time.Sleep(time.Second * 20)

}

func printCacheSize(cache *mCache.Cache) (int, error) {
	return fmt.Println("CacheSize:", cache.Size())
}
