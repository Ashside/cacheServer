package cache

// #include "rocksdb/c.h"
// #include <stdlib.h>
// #cgo CFLAGS: -I${SRCDIR}/rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"errors"
	"regexp"
	"runtime"
	"strconv"
	"unsafe"
)

type RocksdbCache struct {
	db *C.rocksdb_t // RocksDB 存储实例
	ro *C.rocksdb_readoptions_t
	wo *C.rocksdb_writeoptions_t
	e  *C.char // 错误信息
}

// NewRocksdbCache 创建一个RocksdbCache实例
func NewRocksdbCache() *RocksdbCache {
	options := C.rocksdb_options_create()
	// 启动多线程
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU())) // 注意go.int类型转换为C.int
	// 如果不存在则创建
	C.rocksdb_options_set_create_if_missing(options, 1)

	var e *C.char

	// 打开数据库
	db := C.rocksdb_open(options, C.CString("/mnt/rocksdb"), &e)

	if e != nil {
		panic(C.GoString(e))
	}

	C.rocksdb_options_destroy(options)

	return &RocksdbCache{
		db: db,
		ro: C.rocksdb_readoptions_create(),
		wo: C.rocksdb_writeoptions_create(),
		e:  e,
	}
}

// Get 获取键值对
func (r *RocksdbCache) Get(key string) ([]byte, error) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))

	var valLen C.size_t
	v := C.rocksdb_get(r.db, r.ro, k, C.size_t(len(key)), &valLen, &r.e)
	if r.e != nil {
		return nil, errors.New(C.GoString(r.e))
	}
	if v == nil {
		return nil, nil
	}
	defer C.rocksdb_free(unsafe.Pointer(v))

	return C.GoBytes(unsafe.Pointer(v), C.int(valLen)), nil
}

// Set 设置键值对
func (r *RocksdbCache) Set(key string, value []byte) error {
	k := C.CString(key)
	v := C.CBytes(value)
	defer C.free(unsafe.Pointer(k))
	defer C.free(v)

	C.rocksdb_put(r.db, r.wo, k, C.size_t(len(key)), (*C.char)(v), C.size_t(len(value)), &r.e)
	if r.e != nil {
		return errors.New(C.GoString(r.e))
	}
	return nil
}

// Del 删除键值对
func (r *RocksdbCache) Del(key string) error {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))

	C.rocksdb_delete(r.db, r.wo, k, C.size_t(len(key)), &r.e)
	if r.e != nil {
		return errors.New(C.GoString(r.e))
	}
	return nil
}

// GetStat 获取统计信息
func (r *RocksdbCache) GetStat() Stat {
	// 获取rocksdb.aggregated-table-properties属性
	k := C.CString("rocksdb.aggregated-table-properties")
	defer C.free(unsafe.Pointer(k))

	v := C.rocksdb_property_value(r.db, k)
	defer C.free(unsafe.Pointer(v))

	p := C.GoString(v)
	reg := regexp.MustCompile(`([^;]+) = ([^;]+);`) // 表示匹配所有的键值对
	st := Stat{}

	for _, kv := range reg.FindAllStringSubmatch(p, -1) {
		if kv[1] == " # entries" {
			st.Count, _ = strconv.ParseInt(kv[2], 10, 64)
		} else if kv[1] == " raw key size" {
			st.KeySize, _ = strconv.ParseInt(kv[2], 10, 64)
		} else if kv[1] == " raw value size" {
			st.ValueSize, _ = strconv.ParseInt(kv[2], 10, 64)
		}
	}
	return st
}
