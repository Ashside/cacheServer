package cache

type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

func (s *Stat) add(k string, v []byte) {
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
	s.Count++
}

func (s *Stat) del(k string, v []byte) {
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
	s.Count--
}
