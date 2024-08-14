package cache

import "sync"

type inMemoryCache struct {
	cMap  map[string][]byte
	mutex sync.RWMutex
	Stat  // 嵌入Stat结构体
}

func (c *inMemoryCache) Set(k string, v []byte) error {
	// lock
	// 同时只能有一个goroutine访问
	c.mutex.Lock()
	defer c.mutex.Unlock()
	tmp, ok := c.cMap[k]
	if ok {
		// delete the old key
		c.del(k, tmp)
	}
	c.cMap[k] = v
	c.add(k, v)
	return nil
}

func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.cMap[k], nil
}

func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.cMap[k]; ok {
		delete(c.cMap, k)
		c.del(k, v)
		// 两次del操作，第一次是删除map中的key，第二次是删除统计信息
	}
	return nil
}

func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

func newInMemoryCache() Cache {
	c := &inMemoryCache{make(map[string][]byte), sync.RWMutex{}, Stat{}}
	return c
}
