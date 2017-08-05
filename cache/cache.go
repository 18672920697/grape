package cache

import (
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/protocol"
	"sync"
	"strings"
	"fmt"
)

type Cache struct {
	storage *map[string]string
	sync.Mutex
}
func NewCache(c *config.Config) *Cache {
	storage := make(map[string]string)

	cache := &Cache{
		storage: &storage,
	}
	return cache
}

func (cache *Cache) HandleCommand(data protocol.CommandData) {
	switch strings.ToUpper(data.Command) {
	case "SET":
		{

		}
	case "GET":
		{

		}
	}
}

func (cache *Cache) Set(key string, value string) error {
	cache.Lock()
	defer cache.Unlock()
	(*cache.storage)[key] = value

	return nil
}

func (cache *Cache) Get(key string) (string, error) {
	if value, ok := (*cache.storage)[key]; ok {
		return value, nil
	} else {
		return "", fmt.Errorf("Key not found in cache")
	}
}