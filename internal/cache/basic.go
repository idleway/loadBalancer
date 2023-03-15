package cache

import (
	"bytes"
	"sync"
)

const CacheKeyCustomHeader = "LoadBalancer-Cache-Key"

type StoreInCache struct {
	Key      uint64
	RespBody *bytes.Buffer
}

type cache struct {
	data sync.Map
}

func NewCache() *cache {
	obj := cache{
		data: sync.Map{},
	}
	return &obj
}

func (obj *cache) Describe(key uint64) (value []byte, exist bool) {
	valueTemp, exist := obj.data.Load(key)
	if exist {
		value = valueTemp.([]byte)
	}
	return
}
func (obj *cache) Store(key uint64, value *bytes.Buffer) {
	obj.data.Store(key, value.Bytes())
}
