package cache

import "log"

type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() Stat
}

func New(t string) Cache {
	var c Cache
	if t == "inmemory" {
		c = newInMemoryCache()
	}
	if c == nil {
		panic("unknown cache type " + t)
	}
	log.Println(t, "ready to serve")
	return c
}
