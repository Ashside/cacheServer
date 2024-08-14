package http

import (
	"cacheServer/cache"
	"net/http"
)

type Server struct {
	cache.Cache
}

func New(c cache.Cache) *Server {
	return &Server{c}
}

func (s *Server) Listen() {
	// 指定路由和处理函数
	// 前往实现ServeHTTP方法的结构体
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statHandler())
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		return
	}

}
