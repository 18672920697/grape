package cache

import (
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/consistent"
	"github.com/leviathan1995/grape/protocol"

	"fmt"
	"strings"
	"sync"
)

type Cache struct {
	storage     *map[string]string
	config      *config.Config
	consistency *consistent.Consistent
	sync.Mutex
}

func NewCache(config *config.Config, consistency *consistent.Consistent) *Cache {
	storage := make(map[string]string)

	cache := &Cache{
		storage:     &storage,
		config:      config,
		consistency: consistency,
	}
	return cache
}

// Check this key whether store in node
func (cache *Cache) CheckKey(key string) (bool, string) {
	server, _ := cache.consistency.SetKey(key)
	if server != cache.config.Address {
		return false, server
	} else {
		return true, ""
	}
}

func (cache *Cache) HandleCommand(data protocol.CommandData) (protocol.Status, string) {
	switch strings.ToUpper(data.Args[0]) {
	case "COMMAND":
		return protocol.ProtocolNotSupport, ""
	case "SET":
		return cache.HandleSet(data.Args)
	case "GET":
		return cache.HandleGet(data.Args)
	default:
		return protocol.ProtocolNotSupport, ""
	}
}

func (cache *Cache) HandleSet(args []string) (protocol.Status, string) {
	key := args[1]
	// Check this key whether store in node
	if store, server := cache.CheckKey(key); !store {
		return protocol.ProtocolOtherNode, server
	}
	value := args[2]

	cache.Lock()
	defer cache.Unlock()
	(*cache.storage)[key] = value

	resp := fmt.Sprintf("+OK\r\n")
	return protocol.RequestFinish, resp
}

func (cache *Cache) HandleGet(args []string) (protocol.Status, string) {
	key := args[1]
	// Check this key whether store in node
	if store, server := cache.CheckKey(key); !store {
		return protocol.ProtocolOtherNode, server
	}

	if value, ok := (*cache.storage)[key]; ok {
		resp := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
		return protocol.RequestFinish, resp
	} else {
		return protocol.RequestNotFound, ""
	}
}
