package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type cacheHandler struct {
	*Server
}

func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	k := strings.Split(r.URL.EscapedPath(), "/")[2]
	if len(k) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		v, err := h.Get(k)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// 如果没有找到
		if len(v) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, err = w.Write(v)
		if err != nil {
			return
		}
		return
	case http.MethodDelete:
		err := h.Del(k)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	case http.MethodPut:
		v, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(v) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.Set(k, v)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}

}

func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
}

type statHandler struct {
	*Server
}

func (h *statHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	v, err := json.Marshal(h.GetStat())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(v)
	if err != nil {
		return
	}
}

func (s *Server) statHandler() http.Handler {
	return &statHandler{s}
}
