package cache

import (
	"sync"
)

type cacheShard struct {
	shardMap map[uint64]uint32
	dataMap  map[string]string
	sync.RWMutex
}

func NewShard() *cacheShard {
	return &cacheShard{
		shardMap: make(map[uint64]uint32),
		dataMap:  make(map[string]string),
	}
}
