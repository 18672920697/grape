package cache

import (
	"sync"
)

type cacheShard struct {
	dataMap map[string]string
	sync.RWMutex
}

func NewShard() *cacheShard {
	return &cacheShard{
		dataMap: make(map[string]string),
	}
}
