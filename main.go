package main

import (
	"cacheServer/cache"
	"cacheServer/http"
	"cacheServer/tcp"
	"flag"
	"log"
)

func main() {
	typ := flag.String("type", "inmemory", "cache type")
	flag.Parse()
	log.Println("cache type:", *typ)
	c := cache.New(*typ)

	go tcp.New(c).Listen()
	http.New(c).Listen()
}
