package main

import (
	"cacheServer/cache"
	"cacheServer/http"
	"cacheServer/tcp"
)

func main() {
	c := cache.New("inmemory")
	go tcp.New(c).Listen()
	http.New(c).Listen()
}
