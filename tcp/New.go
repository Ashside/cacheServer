package tcp

import "cacheServer/cache"

func New(c cache.Cache) *Server {
	return &Server{c}
}
