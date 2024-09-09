package tcp

import (
	"bufio"
	"cacheServer/cache"
	"io"
	"log"
	"net"
)

type Server struct {
	cache.Cache
}

func (s *Server) Listen() {
	l, e := net.Listen("tcp", ":12346")
	if e != nil {
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e != nil {
			panic(e)
		}
		go s.process(c)
	}
}

func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	log.Println("get")

	key, err := s.readKey(r)
	if err != nil {
		return err
	}
	value, err := s.Get(key)
	log.Println("get key:", key)
	log.Println("get value:", value)

	return sendResponse(value, err, conn)

}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	log.Println("set")
	key, value, err := s.readKeyAndValue(r)
	if err != nil {
		return err

	}
	err = s.Set(key, value)

	log.Println("set key:", key)
	log.Println("set value:", value)

	return sendResponse(nil, err, conn)
}

func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	key, err := s.readKey(r)
	if err != nil {
		return err
	}
	err = s.Del(key)
	return sendResponse(nil, err, conn)

}
func (s *Server) process(conn net.Conn) {

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Close connection due to error:", err)
		}
	}(conn)

	// 使用缓冲区读取数据，避免丢失数据
	// 当我们从conn中读取数据时，可能只读取到了部分数据，剩下的数据还没有读取到，
	// 这时候我们需要使用bufio.Reader来阻塞等待数据的到来，直到读取到了完整的数据
	r := bufio.NewReader(conn)
	for {
		op, err := r.ReadByte()
		if err != nil {
			if err != io.EOF {
				log.Println("Close connection due to error:", err)
				return
			}
		}
		// 根据不同的操作类型，执行不同的操作
		if op == 'S' {
			err = s.set(conn, r)
		} else if op == 'G' {
			err = s.get(conn, r)
		} else if op == 'D' {
			err = s.del(conn, r)
		}
		if err != nil && err != io.EOF {
			log.Println("Close connection due to error:", err)
			return
		}
	}
}
