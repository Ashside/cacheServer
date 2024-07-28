package main

import (
	"net/http"
)

func main() {
	c := cache.New("inmemory")
	http.New(c).Listen()
}
