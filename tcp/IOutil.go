package tcp

import "bufio"

// readKey 读取key
func (s *Server) readKey(r *bufio.Reader) (string, error) {
	// 读取key的长度
	keyLen, e := readLen(r)
	if e != nil {
		return "", e

	}
	key := make([]byte, keyLen)
	// 读取key
	_, e = r.Read(key)

	if e != nil {
		return "", e

	}
	return string(key), nil
}

// readKeyAndValue 读取key和value
func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	keyLen, e := readLen(r)
	if e != nil {
		return "", nil, e
	}
	valueLen, e := readLen(r)
	if e != nil {
		return "", nil, e
	}
	key := make([]byte, keyLen)
	_, e = r.Read(key)
	if e != nil {
		return "", nil, e
	}
	value := make([]byte, valueLen)
	_, e = r.Read(value)
	if e != nil {
		return "", nil, e
	}
	return string(key), value, nil
}
