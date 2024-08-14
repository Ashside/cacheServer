package main

import (
	"cacheServer/cache"
	"cacheServer/http"
)

func main() {
	c := cache.New("inmemory")
	s := http.New(c)
	s.Listen()
}
